/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package portforward

import (
	"context"
	"fmt"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"k8s.io/apimachinery/pkg/labels"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"

	internalLog "github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/store"
)

//go:generate mockgen -source=service.go -destination=./fake/mock_interface.go -package=fake github.com/vmware-tanzu/octant/internal/portforward PortForwarder

var (
	emptyPortForwardResponse = CreateResponse{}
)

// PortForwarder allows querying active port-forwards
type PortForwarder interface {
	List(ctx context.Context) []State
	Get(id string) (State, bool)
	Create(ctx context.Context, gvk schema.GroupVersionKind, name string, namespace string, remotePort uint16) (CreateResponse, error)
	FindTarget(namespace string, gvk schema.GroupVersionKind, name string) ([]State, error)
	FindPod(namespace string, gvk schema.GroupVersionKind, name string) ([]State, error)
	Stop()
	StopForwarder(id string)
}

// PortForwardPortSpec describes a forwarded port.
type PortForwardPortSpec struct {
	Remote uint16 `json:"remote"`
	Local  uint16 `json:"local,omitempty"`
}

// PortForwardSpec describes a port forward.
// TODO Merge with PortForwardState (GH#498)
type PortForwardSpec struct {
	ID        string                `json:"id"`
	Status    string                `json:"status"`
	Message   string                `json:"message"`
	Ports     []PortForwardPortSpec `json:"ports"`
	CreatedAt time.Time             `json:"createdAt"`
}

type CreateRequest struct {
	Namespace  string                `json:"namespace"`
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Name       string                `json:"name"`
	Ports      []PortForwardPortSpec `json:"ports"`
}

type CreateResponse PortForwardSpec

// Target references a kubernetes object
type Target struct {
	GVK       schema.GroupVersionKind
	Namespace string
	Name      string
}

// State describes a single port-forward's runtime state
type State struct {
	ID        string
	CreatedAt time.Time
	Ports     []ForwardedPort
	Target    Target
	Pod       Target

	cancel context.CancelFunc
	ctx    context.Context
}

// Clone clones a port forward state.
func (pf *State) Clone() State {
	pfCpy := State{
		ID:        pf.ID,
		CreatedAt: pf.CreatedAt,
		Ports:     make([]ForwardedPort, len(pf.Ports)),
		Target:    pf.Target,
		Pod:       pf.Pod,
		cancel:    pf.cancel,
		ctx:       pf.ctx,
	}
	copy(pfCpy.Ports, pf.Ports)
	return pfCpy
}

// States describes all active port-forwards' runtime state
type States struct {
	sync.Mutex
	portForwards map[string]State
}

// ServiceOptions contains all the options for running a port-forward service
type ServiceOptions struct {
	RESTClient    rest.Interface
	Config        *restclient.Config
	ObjectStore   store.Store
	PortForwarder portForwarder
}

type forwarderEvent struct {
	ID  string
	err error
}

// Service is a port forwarding service.
type Service struct {
	logger   log.Logger
	opts     ServiceOptions
	ctx      context.Context
	cancel   context.CancelFunc
	notifyCh chan forwarderEvent
	state    States
}

// Check that struct satisfies interface
var _ PortForwarder = (*Service)(nil)

// New creates an instance of Service.
func New(ctx context.Context, opts ServiceOptions) *Service {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		logger:   internalLog.From(ctx),
		opts:     opts,
		notifyCh: make(chan forwarderEvent, 32),
		ctx:      ctx,
		cancel:   cancel,
		state: States{
			portForwards: make(map[string]State),
		},
	}
}

// Stop stops all forwarders. The portForwardService is invalid after calling stop.
func (s *Service) Stop() {
	// TODO wait on goroutines to complete after calling cancel. (GH#494)
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *Service) validateCreateRequest(r CreateRequest) error {
	if r.Namespace == "" {
		return errors.New("namespace field required")
	}
	if r.Name == "" {
		return errors.New("name field required")
	}

	if r.APIVersion != "v1" || (r.Kind != "Pod" && r.Kind != "Service") {
		return errors.Errorf("port forwards only work with pods & services")
	}

	for _, p := range r.Ports {
		if p.Remote < 1 || p.Remote > 65535 {
			return errors.Errorf("remote port out of range: %v", p.Remote)
		}
	}

	return nil
}

// resolvePod attempts to resolve a port forward request into an active pod we can
// forward to. Service/deployments selectors will be resolved into pods and a random
// one will be chosen. A pod has to be active.
// Returns: pod name or error.
func (s *Service) resolvePod(ctx context.Context, r CreateRequest) (string, error) {
	o := s.opts.ObjectStore
	if o == nil {
		return "", errors.New("nil objectstore")
	}

	switch {
	case r.APIVersion == "v1" && r.Kind == "Pod":
		// Verify pod exists and status is running
		if ok, err := s.verifyPod(ctx, r.Namespace, r.Name); !ok || err != nil {
			return "", errors.Errorf("verifying pod %q: %v", r.Name, err)
		}
		return r.Name, nil
	case r.APIVersion == "v1" && r.Kind == "Service":
		pod, err := s.findPodForService(ctx, r.APIVersion, r.Kind, r.Namespace, r.Name)
		if err != nil {
			return "", err
		}

		return pod.Name, nil
	default:
		return "", errors.New("not implemented")
	}

}

func (s *Service) findPodForService(ctx context.Context, apiVersion, kind, namespace, name string) (*corev1.Pod, error) {
	o := s.opts.ObjectStore
	if o == nil {
		return nil, errors.New("nil objectstore")
	}

	key := store.Key{
		APIVersion: apiVersion,
		Kind:       kind,
		Namespace:  namespace,
		Name:       name,
	}
	var service corev1.Service
	found, err := store.GetAs(ctx, o, key, &service)

	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	lbls := labels.Set(service.Spec.Selector)
	key = store.Key{
		APIVersion: apiVersion,
		Kind:       "Pod",
		Namespace:  namespace,
		Selector:   &lbls,
	}
	list, _, err := o.List(ctx, key)
	if err != nil {
		return nil, err
	}

	for i := range list.Items {
		pod := &corev1.Pod{}

		if err := kubernetes.FromUnstructured(&list.Items[i], pod); err != nil {
			return nil, err
		}

		return pod, nil
	}

	return nil, errors.New("no matching pod found for service")
}

// verifyPod returns true if the specified pod can be found and is in the running phase.
// Otherwise returns false and an error describing the cause.
func (s *Service) verifyPod(ctx context.Context, namespace, name string) (bool, error) {
	o := s.opts.ObjectStore
	if o == nil {
		return false, errors.New("nil objectstore")
	}

	key := store.Key{
		APIVersion: "v1",
		Kind:       "Pod",
		Namespace:  namespace,
		Name:       name,
	}
	var pod corev1.Pod
	found, err := store.GetAs(ctx, o, key, &pod)
	if err != nil {
		return false, err
	}
	if !found {
		return false, nil
	}

	if pod.Name == "" {
		return false, errors.New("pod not found")
	}

	if pod.Status.Phase != corev1.PodRunning {
		return false, errors.Errorf("pod not running, phase=%v", pod.Status.Phase)
	}

	return true, nil
}

// createForwarder creates a port forwarder, forwards traffic, and blocks until
// port state information is populated.
// Returns forwarder id.
func (s *Service) createForwarder(targetRequest, podRequest CreateRequest) (string, error) {
	logger := s.logger.With("context", "PortForwardService.createForwarder")

	if s.opts.PortForwarder == nil {
		return "", errors.New("portforwarder is nil")
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating uuid")
	}
	forwarderID := randomUUID.String()
	logger = logger.With("id", forwarderID)

	var ports []string
	for _, p := range podRequest.Ports {
		ports = append(ports, fmt.Sprintf("%d:%d", p.Local, p.Remote))
	}

	// Target coordinates to preserve in state
	targetGv, err := schema.ParseGroupVersion(targetRequest.APIVersion)
	if err != nil {
		return "", errors.Wrap(err, "parsing APIVersion")
	}
	targetGvk := targetGv.WithKind(targetRequest.Kind)

	// Pod coordinates to preserve in state
	podGv, err := schema.ParseGroupVersion(podRequest.APIVersion)
	if err != nil {
		return "", errors.Wrap(err, "parsing APIVersion")
	}
	podGvk := podGv.WithKind(podRequest.Kind)

	// This child context will be cancelled if our parent context is cancelled
	ctx, cancel := context.WithCancel(s.ctx)

	// Spawns goroutine to update state as ports become available
	portsChannel, portsReady := s.localPortsHandler(ctx, forwarderID)

	o := &s.opts
	opts := Options{
		Config:        o.Config,
		RESTClient:    o.RESTClient,
		Address:       []string{"localhost"},
		Ports:         ports,
		PortForwarder: o.PortForwarder,
		StopChannel:   ctx.Done(),
		ReadyChannel:  make(chan struct{}),
		PortsChannel:  portsChannel,
	}

	// NOTE: ports will be updated in the state struct by
	// localPortsHandler when they become available.
	forwardState := State{
		ID:        forwarderID,
		CreatedAt: time.Now(),
		Target: Target{
			GVK:       targetGvk,
			Namespace: targetRequest.Namespace,
			Name:      targetRequest.Name,
		},
		Pod: Target{
			GVK:       podGvk,
			Namespace: podRequest.Namespace,
			Name:      podRequest.Name,
		},

		cancel: cancel,
		ctx:    ctx,
	}

	s.state.Lock()
	s.state.portForwards[forwarderID] = forwardState
	s.state.Unlock()

	req := o.RESTClient.Post().
		Resource("pods").
		Namespace(podRequest.Namespace).
		Name(podRequest.Name).
		SubResource("portforward")

	go func() {
		// Blocks until forwarder completes
		logger.With("url", req.URL()).Debugf("starting port-forward")
		err := s.opts.PortForwarder.ForwardPorts("POST", req.URL(), opts)

		logger.Debugf("forwarding terminated: %v", err)

		// Notify the main forwarder of the termination
		event := forwarderEvent{
			ID:  forwarderID,
			err: err,
		}
		select {
		case s.notifyCh <- event:
		default:
		}

		// Cleanup state for terminated port-forward
		s.StopForwarder(forwarderID)
	}()

	// Block until ports state is ready
	select {
	case <-ctx.Done():
		return "", errors.Errorf("portforward terminated due to parent context: %v", forwarderID)
	case <-portsReady:
	}

	return forwarderID, nil
}

// responseForCreate creates a create response based on the state for the specified forward (by id)
func (s *Service) responseForCreate(id string) (CreateResponse, error) {
	var response CreateResponse

	s.state.Lock()
	defer s.state.Unlock()
	state, ok := s.state.portForwards[id]
	if !ok {
		return response, errors.Errorf("retrieving state for terminated port-forward: %v", id)
	}

	response.ID = id
	response.CreatedAt = state.CreatedAt
	rp := make([]PortForwardPortSpec, len(state.Ports))
	for i := range state.Ports {
		rp[i].Local = state.Ports[i].Local
		rp[i].Remote = state.Ports[i].Remote
	}
	response.Ports = rp
	response.Status = "ok"
	return response, nil
}

func (s *Service) localPortsHandler(ctx context.Context, id string) (portsChan chan []ForwardedPort, portsReady <-chan struct{}) {
	logger := s.logger.With("context", "PortForwardService.localPortsHandler", "id", id)
	portsChan = make(chan []ForwardedPort, 1)
	readyChan := make(chan struct{})
	portsReady = readyChan

	go func() {
		select {
		case p := <-portsChan:
			logger.With("ports", p).Debugf("received ports for port-forward")
			if err := s.updatePorts(id, p); err != nil {
				logger.Warnf("%s", err.Error())
			}

			close(readyChan)

		case <-ctx.Done():
			logger.Debugf("terminated")
		}
	}()

	return
}

// updatePorts updates the ports list for an existing port forward, specified by id
func (s *Service) updatePorts(id string, ports []ForwardedPort) error {
	s.state.Lock()
	defer s.state.Unlock()
	state, ok := s.state.portForwards[id]
	if !ok {
		return errors.New("updating ports for terminated port-forward")
	}
	state.Ports = ports
	s.state.portForwards[id] = state
	return nil
}

// List lists port forwards
func (s *Service) List(ctx context.Context) []State {
	s.state.Lock()
	defer s.state.Unlock()

	result := make([]State, 0, len(s.state.portForwards))
	for i, pf := range s.state.portForwards {
		targetPod := &pf.Pod
		if verified, err := s.verifyPod(ctx, targetPod.Namespace, targetPod.Name); !verified || err != nil {
			delete(s.state.portForwards, i)
			continue
		}
		result = append(result, pf.Clone())
	}

	// Sort by creation time
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result
}

// Get gets a single port forward state
func (s *Service) Get(id string) (State, bool) {
	s.state.Lock()
	defer s.state.Unlock()

	pf, ok := s.state.portForwards[id]
	cpy := pf.Clone()
	return cpy, ok
}

// Create creates a new port forward for the specified object and remote port.
// Implements PortForwardInterface.
func (s *Service) Create(ctx context.Context, gvk schema.GroupVersionKind, name string, namespace string, remotePort uint16) (CreateResponse, error) {
	logger := s.logger.With("context", "PortForwardService.Create")
	req := newForwardRequest(gvk, name, namespace, remotePort)

	if err := s.validateCreateRequest(req); err != nil {
		return emptyPortForwardResponse, errors.Wrap(err, "invalid request")
	}

	// Resolve the request into a pod, update the request
	logger.With(
		"apiVersion", req.APIVersion,
		"kind", req.Kind,
		"name", req.Name,
		"namespace", req.Namespace,
	).Debugf("resolving pod from object")
	podName, err := s.resolvePod(ctx, req)
	if err != nil {
		return emptyPortForwardResponse, errors.Wrap(err, "resolving pod")
	}
	logger.Debugf("resolved to pod %q", podName)
	podReq := req
	podReq.Name = podName
	podReq.Kind = "Pod"

	id, err := s.createForwarder(req, CreateRequest{
		Namespace:  req.Namespace,
		APIVersion: req.APIVersion,
		Kind:       "Pod",
		Name:       podName,
		Ports:      req.Ports,
	})
	if err != nil {
		return emptyPortForwardResponse, errors.Wrap(err, "creating forwarder")
	}

	// Compose response based on forwarder state
	response, err := s.responseForCreate(id)
	if err != nil {
		return emptyPortForwardResponse, errors.Wrapf(err, "fetching state for forwarder: %v", id)
	}

	return response, nil
}

// StopForwarder stops an individual port forward specified by id.
// Implements PortForwardInterface.
func (s *Service) StopForwarder(id string) {
	s.state.Lock()
	defer s.state.Unlock()

	pf, ok := s.state.portForwards[id]
	if !ok {
		return
	}
	if pf.cancel != nil {
		pf.cancel()
		<-pf.ctx.Done() // Wait for PortForward context to finish.
	}

	delete(s.state.portForwards, id)
}

type notFound struct{}

// Check that struct satisfies interface
var _ error = (*notFound)(nil)

func (e *notFound) Error() string {
	return "port forward not found"
}

func (e *notFound) NotFound() bool {
	return true
}

func (s *Service) FindTarget(namespace string, gvk schema.GroupVersionKind, name string) ([]State, error) {
	s.state.Lock()
	defer s.state.Unlock()

	result := make([]State, 0, len(s.state.portForwards))

	for _, state := range s.state.portForwards {
		target := state.Target
		if target.GVK.String() == gvk.String() &&
			namespace == target.Namespace &&
			name == target.Name {
			result = append(result, state)
		}
	}

	return result, &notFound{}
}

func (s *Service) FindPod(namespace string, gvk schema.GroupVersionKind, name string) ([]State, error) {
	s.state.Lock()
	defer s.state.Unlock()

	result := make([]State, 0, len(s.state.portForwards))

	for _, state := range s.state.portForwards {
		target := state.Pod
		if target.GVK.String() == gvk.String() &&
			namespace == target.Namespace &&
			name == target.Name {
			result = append(result, state)
		}
	}

	return result, &notFound{}
}

// newForwardRequest constructs a port forwarding request based on the provided parameters
func newForwardRequest(gvk schema.GroupVersionKind, name string, namespace string, remotePort uint16) CreateRequest {
	APIVersion, kind := gvk.ToAPIVersionAndKind()

	return CreateRequest{
		APIVersion: APIVersion,
		Kind:       kind,
		Namespace:  namespace,
		Name:       name,
		Ports: []PortForwardPortSpec{
			{
				Remote: remotePort,
			},
		},
	}
}
