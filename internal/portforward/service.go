package portforward

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/log"
	"github.com/heptio/developer-dash/internal/mime"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
)

//go:generate mockgen -destination=./fake/mock_interface.go -package=fake github.com/heptio/developer-dash/internal/portforward PortForwardInterface

var (
	emptyPortForwardResponse = PortForwardCreateResponse{}
)

// PortForwardInterface allows querying active port-forwards
type PortForwardInterface interface {
	List() []PortForwardState
	Get(id string) (PortForwardState, bool)
	Create(ctx context.Context, gvk schema.GroupVersionKind, name string, namespace string, remotePort uint16) (PortForwardCreateResponse, error)
	Find(namespace string, gvk schema.GroupVersionKind, name string) (PortForwardState, error)
	Stop()
	StopForwarder(id string)
}

type portForwardPortSpec struct {
	Remote uint16 `json:"remote"`
	Local  uint16 `json:"local,omitempty"`
}

// TODO Merge with PortForwardState
type portForwardSpec struct {
	ID        string                `json:"id"`
	Status    string                `json:"status"`
	Message   string                `json:"message"`
	Ports     []portForwardPortSpec `json:"ports"`
	CreatedAt time.Time             `json:"createdAt"`
}

type PortForwardCreateRequest struct {
	Namespace  string                `json:"namespace"`
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Name       string                `json:"name"`
	Ports      []portForwardPortSpec `json:"ports"`
}

type PortForwardCreateResponse portForwardSpec

type portForwardListRequest struct {
	ID string
}

type portForwardDeleteRequest struct {
	ID string
}

type portForwardListResponse struct {
	PortForwards []portForwardSpec `json:"portforwards"`
}

type portForwardDeleteResponse struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PortForwardTarget references a kubernetes object
type PortForwardTarget struct {
	GVK       schema.GroupVersionKind
	Namespace string
	Name      string
}

// PortForwardState describes a single port-forward's runtime state
type PortForwardState struct {
	ID        string
	CreatedAt time.Time
	Ports     []ForwardedPort
	Target    PortForwardTarget
	Pod       PortForwardTarget

	cancel context.CancelFunc
}

// Clone clones a port forward state.
func (pf *PortForwardState) Clone() PortForwardState {
	pfCpy := PortForwardState{
		ID:        pf.ID,
		CreatedAt: pf.CreatedAt,
		Ports:     make([]ForwardedPort, len(pf.Ports)),
		Target:    pf.Target,
		Pod:       pf.Pod,
		cancel:    pf.cancel,
	}
	copy(pfCpy.Ports, pf.Ports)
	return pfCpy
}

// PortForwardStates describes all active port-forwards' runtime state
type PortForwardStates struct {
	sync.Mutex
	portForwards map[string]PortForwardState
}

// PortForwardSvcOptions contains all the options for running a port-forward service
type PortForwardSvcOptions struct {
	RESTClient    rest.Interface
	Config        *restclient.Config
	Cache         cache.Cache
	PortForwarder portForwarder
}

type forwarderEvent struct {
	ID  string
	err error
}

type PortForwardService struct {
	logger   log.Logger
	opts     PortForwardSvcOptions
	ctx      context.Context
	cancel   context.CancelFunc
	notifyCh chan forwarderEvent
	state    PortForwardStates
}

var _ PortForwardInterface = (*PortForwardService)(nil)

func NewPortForwardService(ctx context.Context, opts PortForwardSvcOptions, logger log.Logger) *PortForwardService {
	ctx, cancel := context.WithCancel(ctx)
	return &PortForwardService{
		logger:   logger,
		opts:     opts,
		notifyCh: make(chan forwarderEvent, 32),
		ctx:      ctx,
		cancel:   cancel,
		state: PortForwardStates{
			portForwards: make(map[string]PortForwardState),
		},
	}
}

// Stop stops all forwarders. The portForwardService is invalid after calling stop.
func (s *PortForwardService) Stop() {
	// TODO wait on goroutines to complete after calling cancel.
	if s.cancel != nil {
		s.cancel()
	}
}

func (s *PortForwardService) validateCreateRequest(r PortForwardCreateRequest) error {
	if r.Namespace == "" {
		return errors.New("namespace field required")
	}
	if r.Name == "" {
		return errors.New("name field required")
	}
	switch r.APIVersion {
	case "v1":
	case "apps/v1":
	default:
		return fmt.Errorf("unsupported apiVersion (%s) - must be one of (v1, apps/v1)", r.APIVersion)
	}

	switch r.Kind {
	case "Deployment":
		if r.APIVersion != "v1" {
			return fmt.Errorf("Unsupported resource: %s/%s", r.APIVersion, r.Kind)
		}
	case "Service":
		if r.APIVersion != "apps/v1" {
			return fmt.Errorf("Unsupported resource: %s/%s", r.APIVersion, r.Kind)
		}
	case "Pod":
		if r.APIVersion != "v1" {
			return fmt.Errorf("Unsupported resource: %s/%s", r.APIVersion, r.Kind)
		}
	default:
		return fmt.Errorf("unsupported kind (%s) - must be one of (Deployment, Service, Pod)", r.Kind)
	}

	for _, p := range r.Ports {
		if p.Remote < 1 || p.Remote > 65535 {
			return fmt.Errorf("remote port out of range: %v", p.Remote)
		}
	}

	return nil
}

// resolvePod attempts to resolve a port forward request into an active pod we can
// forward to. Service/deployments selectors will be resolved into pods and a random
// one will be chosen. A pod has to be active.
// Returns: pod name or error.
func (s *PortForwardService) resolvePod(ctx context.Context, r PortForwardCreateRequest) (string, error) {
	c := s.opts.Cache
	if c == nil {
		return "", errors.New("nil cache")
	}

	switch {
	case r.APIVersion == "v1" && r.Kind == "Pod":
		// Verify pod exists and status is running
		if ok, err := s.verifyPod(ctx, r.Namespace, r.Name); !ok || err != nil {
			return "", fmt.Errorf("verifying pod %q: %v", r.Name, err)
		}
		return r.Name, nil
	case r.APIVersion == "v1" && r.Kind == "Service":
		// TODO: implement service, deployment cases
		return "", errors.New("not implemented")
	case r.APIVersion == "apps/v1" && r.Kind == "Deployment":
		// TODO: implement service, deployment cases
		return "", errors.New("not implemented")
	default:
		return "", errors.New("not implemented")
	}

}

// verifyPod returns true if the specified pod can be found and is in the running phase.
// Otherwise returns false and an error describing the cause.
func (s *PortForwardService) verifyPod(ctx context.Context, namespace, name string) (bool, error) {
	c := s.opts.Cache
	if c == nil {
		return false, errors.New("nil cache")
	}

	key := cache.Key{
		APIVersion: "v1",
		Kind:       "Pod",
		Namespace:  namespace,
		Name:       name,
	}
	var pod corev1.Pod
	if err := cache.GetAs(ctx, c, key, &pod); err != nil {
		return false, err
	}
	if pod.Name == "" {
		return false, errors.New("pod not found")
	}

	if pod.Status.Phase != corev1.PodRunning {
		return false, fmt.Errorf("pod not running, phase=%v", pod.Status.Phase)
	}

	return true, nil
}

// createForwarder creates a port forwarder, forwards traffic, and blocks until
// port state information is populated.
// Returns forwarder id.
func (s *PortForwardService) createForwarder(r PortForwardCreateRequest) (string, error) {
	log := s.logger.With("context", "PortForwardService.createForwarder")

	if s.opts.PortForwarder == nil {
		return "", errors.New("portforwarder is nil")
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "generating uuid")
	}
	idstr := id.String()
	log = log.With("id", idstr)

	var ports []string
	for _, p := range r.Ports {
		ports = append(ports, fmt.Sprintf("%d:%d", p.Local, p.Remote))
	}

	// Target coordinates to preserve in state
	gv, err := schema.ParseGroupVersion(r.APIVersion)
	if err != nil {
		return "", errors.Wrap(err, "parsing APIVersion")
	}
	gvk := gv.WithKind(r.Kind)

	// This child context will be cancelled if our parent context is cancelled
	ctx, cancel := context.WithCancel(s.ctx)

	// Spawns goroutine to update state as ports become available
	portsChannel, portsReady := s.localPortsHandler(ctx, idstr)

	// TODO resolve request gvk/name to pod name
	o := &s.opts
	opts := PortForwardOptions{
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
	forwardState := PortForwardState{
		ID:        idstr,
		CreatedAt: time.Now(),
		Target: PortForwardTarget{
			GVK:       gvk,
			Namespace: r.Namespace,
			Name:      r.Name,
		},
		// TODO Target and Pod may be different
		Pod: PortForwardTarget{
			GVK:       gvk,
			Namespace: r.Namespace,
			Name:      r.Name,
		},
		cancel: cancel,
	}

	s.state.Lock()
	s.state.portForwards[idstr] = forwardState
	s.state.Unlock()

	req := o.RESTClient.Post().
		Resource("pods").
		Namespace(r.Namespace).
		Name(r.Name).
		SubResource("portforward")

	go func() {
		// Blocks until forwarder completes
		log.With("url", req.URL()).Debugf("starting port-forward")
		err := s.opts.PortForwarder.ForwardPorts("POST", req.URL(), opts)

		log.Debugf("forwarding terminated: %v", err)

		// Notify the main forwarder of the termination
		event := forwarderEvent{
			ID:  idstr,
			err: err,
		}
		select {
		case s.notifyCh <- event:
		default:
		}

		// Cleanup state for terminated port-forward
		s.StopForwarder(idstr)
	}()

	// Block until ports state is ready
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("portforward terminated due to parent context: %v", idstr)
	case <-portsReady:
	}

	return idstr, nil
}

// responseForCreate creates a create response based on the state for the specified forward (by id)
func (s *PortForwardService) responseForCreate(id string) (PortForwardCreateResponse, error) {
	var response PortForwardCreateResponse

	s.state.Lock()
	defer s.state.Unlock()
	state, ok := s.state.portForwards[id]
	if !ok {
		return response, fmt.Errorf("retrieving state for terminated port-forward: %v", id)
	}

	response.ID = id
	response.CreatedAt = state.CreatedAt
	rp := make([]portForwardPortSpec, len(state.Ports))
	for i := range state.Ports {
		rp[i].Local = state.Ports[i].Local
		rp[i].Remote = state.Ports[i].Remote
	}
	response.Ports = rp
	response.Status = "ok"
	return response, nil
}

func (s *PortForwardService) localPortsHandler(ctx context.Context, id string) (portsChan chan []ForwardedPort, portsReady <-chan struct{}) {
	log := s.logger.With("context", "PortForwardService.localPortsHandler", "id", id)
	portsChan = make(chan []ForwardedPort, 1)
	readyChan := make(chan struct{})
	portsReady = readyChan

	go func() {
		select {
		case p := <-portsChan:
			log.With("ports", p).Debugf("received ports for port-forward")
			if err := s.updatePorts(id, p); err != nil {
				log.Warnf("%s", err.Error())
			}

			close(readyChan)

		case <-ctx.Done():
			log.Debugf("terminated")
		}
	}()

	return
}

// updatePorts updates the ports list for an existing port forward, specified by id
func (s *PortForwardService) updatePorts(id string, ports []ForwardedPort) error {
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

func (s *PortForwardService) stateByID(id string) (PortForwardState, bool) {
	s.state.Lock()
	defer s.state.Unlock()
	pf, ok := s.state.portForwards[id]
	return pf, ok
}

// List lists port forwards
func (s *PortForwardService) List() []PortForwardState {
	s.state.Lock()
	defer s.state.Unlock()

	result := make([]PortForwardState, 0, len(s.state.portForwards))
	for _, pf := range s.state.portForwards {
		result = append(result, pf.Clone())
	}

	// Sort by creation time
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result
}

// Get gets a single port forward state
func (s *PortForwardService) Get(id string) (PortForwardState, bool) {
	s.state.Lock()
	defer s.state.Unlock()

	pf, ok := s.state.portForwards[id]
	cpy := pf.Clone()
	return cpy, ok
}

// Create creates a new port forward for the specified object and remote port.
// Implements PortForwardInterface.
func (s *PortForwardService) Create(ctx context.Context, gvk schema.GroupVersionKind, name string, namespace string, remotePort uint16) (PortForwardCreateResponse, error) {
	log := s.logger.With("context", "PortForwardService.Create")
	req := newForwardRequest(gvk, name, namespace, remotePort)

	if err := s.validateCreateRequest(req); err != nil {
		return emptyPortForwardResponse, errors.Wrap(err, "invalid request")
	}

	// Resolve the request into a pod, update the request
	log.With(
		"apiVersion", req.APIVersion,
		"kind", req.Kind,
		"name", req.Name,
		"namespace", req.Namespace,
	).Debugf("resolving pod from object")
	podName, err := s.resolvePod(ctx, req)
	if err != nil {
		return emptyPortForwardResponse, errors.Wrap(err, "resolving pod")
	}
	log.Debugf("resolved to pod %q", podName)
	podReq := req
	podReq.Name = podName

	id, err := s.createForwarder(req)
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
func (s *PortForwardService) StopForwarder(id string) {
	s.state.Lock()
	defer s.state.Unlock()

	pf, ok := s.state.portForwards[id]
	if !ok {
		return
	}
	if pf.cancel != nil {
		// TODO wait for goroutine to exit
		pf.cancel()
	}
	delete(s.state.portForwards, id)
}

type notFound struct{}

var _ error = (*notFound)(nil)

func (e *notFound) Error() string {
	return "port forward not found"
}

func (e *notFound) NotFound() bool {
	return true
}

func (s *PortForwardService) Find(namespace string, gvk schema.GroupVersionKind, name string) (PortForwardState, error) {
	s.state.Lock()
	defer s.state.Unlock()

	for _, state := range s.state.portForwards {
		target := state.Target
		if target.GVK.String() == gvk.String() &&
			namespace == target.Namespace &&
			name == target.Name {
			return state, nil
		}
	}

	return PortForwardState{}, &notFound{}
}

type errorMessage struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type errorResponse struct {
	Error errorMessage `json:"error,omitempty"`
}

// respondWithError - same as api.responsdWithError
func respondWithError(w http.ResponseWriter, code int, message string, logger log.Logger) {
	r := &errorResponse{
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", mime.JSONContentType)

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(r); err != nil {
		logger.Errorf("encoding JSON response: %v", err)
	}
}

// newForwardRequest constructs a port forwarding request based on the provided parameters
func newForwardRequest(gvk schema.GroupVersionKind, name string, namespace string, remotePort uint16) PortForwardCreateRequest {
	APIVersion, kind := gvk.ToAPIVersionAndKind()

	return PortForwardCreateRequest{
		APIVersion: APIVersion,
		Kind:       kind,
		Namespace:  namespace,
		Name:       name,
		Ports: []portForwardPortSpec{
			portForwardPortSpec{
				Remote: remotePort,
			},
		},
	}
}
