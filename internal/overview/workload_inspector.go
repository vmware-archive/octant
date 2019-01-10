package overview

import (
	"context"
	"fmt"
	"github.com/heptio/developer-dash/internal/cache"
	"reflect"
	"sort"
	"strings"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
)

type visitFunc func(context.Context, runtime.Object, cache.Cache, content.Nodes, content.AdjList, visitSet) error

// resourceVisitor visits resources by type and function.
type resourceVisitor struct {
	visitors map[runtime.Object]visitFunc
}

func (rv *resourceVisitor) visit(ctx context.Context, object runtime.Object, c cache.Cache, visited visitSet) ([]content.Content, error) {
	dag := content.NewDAG()
	acc := meta.NewAccessor()

	for visitor, fn := range rv.visitors {
		if reflect.TypeOf(visitor) != reflect.TypeOf(object) {
			continue
		}

		uid, err := acc.UID(object)
		if err != nil {
			return nil, errors.Wrapf(err, "fetching UID for object type %T", object)
		}

		dag.Selected = string(uid)

		if err := fn(ctx, object, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}

		return []content.Content{dag}, nil
	}

	return nil, errors.Errorf("unable to visit resource of type %T", object)
}

var (
	defaultResourceVisitor = &resourceVisitor{
		visitors: map[runtime.Object]visitFunc{
			&corev1.Pod{}:                   visitPodRoot,
			&corev1.Service{}:               visitService,
			&appsv1.Deployment{}:            visitDeployment,
			&appsv1.ReplicaSet{}:            visitReplicaSet,
			&v1beta1.Ingress{}:              visitIngress,
			&appsv1.StatefulSet{}:           visitStatefulSet,
			&corev1.ReplicationController{}: visitReplicationController,
			&appsv1.DaemonSet{}:             visitDaemonSet,
		},
	}
)

type workloadChecks = map[runtime.Object]*nodeStatus

var (
	defaultChecks = workloadChecks{
		&appsv1.Deployment{}: newNodeStatus(deploymentCheckUnavailable),
		&appsv1.ReplicaSet{}: newNodeStatus(replicasSetCheckAvailableReplicas),
	}
)

// WorkloadInspector is both a View and a View
type workloadInspectorView struct {
	workloadChecks workloadChecks
}

type visitKey struct {
	uid k8stypes.UID
}
type visitSet map[visitKey]bool

type workloadInspectorViewOpt func(*workloadInspectorView)

func setViewChecks(m workloadChecks) workloadInspectorViewOpt {
	return func(wiv *workloadInspectorView) {
		wiv.workloadChecks = m
	}
}

// workloadViewFactory creates a view for workload inspect view.
func workloadViewFactory(prefix, namespace string, c clock.Clock) View {
	return newWorkloadInspectorView(prefix, namespace, c)
}

// newWorkloadInspectorView creates an instance of workloadInspectorView
func newWorkloadInspectorView(prefix, namespace string, c clock.Clock, opts ...workloadInspectorViewOpt) *workloadInspectorView {
	wiv := &workloadInspectorView{
		workloadChecks: defaultChecks,
	}

	for _, opt := range opts {
		opt(wiv)
	}

	return wiv
}

func (wiv *workloadInspectorView) checkResource(r runtime.Object) (ResourceStatusList, error) {
	for k, v := range wiv.workloadChecks {
		if reflect.TypeOf(r) == reflect.TypeOf(k) {
			return v.check(r)
		}
	}

	return ResourceStatusList{}, nil
}

// Implements View.Content
func (wiv *workloadInspectorView) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	visited := visitSet{}

	// TODO: pass the resource visitor in instead of assuming the default will always be used.
	return defaultResourceVisitor.visit(ctx, object, c, visited)
}

func visitKeyForObject(obj runtime.Object) visitKey {
	acc := meta.NewAccessor()
	// TODO ERROR HANDLING
	uid, _ := acc.UID(obj)
	return visitKey{uid}
}

// An edgeFunc will create an edge to the provided destination node
type edgeFunc func(dst string)

func listIngressPaths(ingress *v1beta1.Ingress, c cache.Cache) ([]v1beta1.HTTPIngressPath, error) {
	if ingress == nil {
		return nil, errors.New("nil ingress")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var paths []v1beta1.HTTPIngressPath

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}
		for _, p := range rule.IngressRuleValue.HTTP.Paths {
			paths = append(paths, p)
		}
	}

	return paths, nil
}
func listIngressBackends(ingress *v1beta1.Ingress, c cache.Cache) ([]v1beta1.IngressBackend, error) {
	if ingress == nil {
		return nil, errors.New("nil ingress")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var backends []v1beta1.IngressBackend

	if ingress.Spec.Backend != nil && ingress.Spec.Backend.ServiceName != "" {
		backends = append(backends, *ingress.Spec.Backend)
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}
		for _, p := range rule.IngressRuleValue.HTTP.Paths {
			if p.Backend.ServiceName == "" {
				continue
			}
			backends = append(backends, p.Backend)
		}
	}

	return backends, nil
}

func loadServices(serviceNames []string, namespace string, c cache.Cache) ([]*corev1.Service, error) {
	var services []*corev1.Service
	for _, backend := range serviceNames {
		key := cache.CacheKey{
			Namespace:  namespace,
			APIVersion: "v1",
			Kind:       "Service",
			Name:       backend,
		}
		ul, err := c.Retrieve(key)
		if err != nil {
			return nil, errors.Wrapf(err, "retrieving service backend: %v", backend)
		}
		for _, u := range ul {
			svc := &corev1.Service{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, svc)
			if err != nil {
				return nil, errors.Wrap(err, "converting unstructured service")
			}
			if err := copyObjectMeta(svc, u); err != nil {
				return nil, errors.Wrap(err, "copying object metadata")
			}
			services = append(services, svc)
		}
	}
	return services, nil
}

func loadService(name string, namespace string, c cache.Cache) (*corev1.Service, error) {
	services, err := loadServices([]string{name}, namespace, c)
	if err != nil {
		return nil, err
	}
	if len(services) < 1 {
		return nil, nil
	}
	return services[0], nil
}

// Reverse-lookup ingresses that point to a service
func findIngressesForService(svc *corev1.Service, c cache.Cache) ([]*v1beta1.Ingress, error) {
	var results []*v1beta1.Ingress
	if svc == nil {
		return nil, errors.New("nil service")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := cache.CacheKey{
		Namespace:  svc.Namespace,
		APIVersion: "extensions/v1beta1",
		Kind:       "Ingress",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving ingresses")
	}
	for _, u := range ul {
		ingress := &v1beta1.Ingress{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, ingress)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured ingress")
		}
		if err = copyObjectMeta(ingress, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		backends, err := listIngressBackends(ingress, c)
		if err != nil {
			return nil, errors.Wrapf(err, "listing backends for ingress: %v", ingress.Name)
		}
		if !containsBackend(backends, svc.Name) {
			continue
		}

		results = append(results, ingress)
	}
	return results, nil
}

func findPodsForService(svc *corev1.Service, c cache.Cache) ([]*corev1.Pod, error) {
	if svc == nil {
		return nil, errors.New("nil service")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}
	key := cache.CacheKey{
		Namespace:  svc.Namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	selector, err := getSelector(svc)
	if err != nil {
		return nil, errors.Wrapf(err, "creating pod selector for service: %v", svc.Name)
	}
	pods, err := loadPods(key, c, selector)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching pods for service: %v", svc.Name)
	}

	return pods, nil
}

// Reverse-lookup services that point to a pod
func findServicesForPod(pod *corev1.Pod, c cache.Cache) ([]*corev1.Service, error) {
	var results []*corev1.Service
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := cache.CacheKey{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "Service",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving services")
	}
	for _, u := range ul {
		svc := &corev1.Service{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, svc)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured service")
		}
		if err = copyObjectMeta(svc, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		labelSelector, err := getSelector(svc)
		if err != nil {
			return nil, errors.Wrapf(err, "creating pod selector for service: %v", svc.Name)
		}
		selector, err := metav1.LabelSelectorAsSelector(labelSelector)
		if err != nil {
			return nil, errors.Wrap(err, "invalid selector")
		}

		if selector.Empty() || !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}
		results = append(results, svc)
	}
	return results, nil
}

// Reverse-lookup replicasets that point to a pod
func findReplicaSetsForPod(pod *corev1.Pod, c cache.Cache) ([]*appsv1.ReplicaSet, error) {
	var results []*appsv1.ReplicaSet
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := cache.CacheKey{
		Namespace:  pod.Namespace,
		APIVersion: "apps/v1",
		Kind:       "ReplicaSet",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving replicaSets")
	}
	for _, u := range ul {
		rs := &appsv1.ReplicaSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, rs)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured replicaSet")
		}
		if err = copyObjectMeta(rs, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		labelSelector, err := getSelector(rs)
		if err != nil {
			return nil, errors.Wrapf(err, "creating pod selector for replicaSet: %v", rs.Name)
		}
		selector, err := metav1.LabelSelectorAsSelector(labelSelector)
		if err != nil {
			return nil, errors.Wrap(err, "invalid selector")
		}

		if selector.Empty() || !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}
		results = append(results, rs)
	}
	return results, nil
}

// Reverse-lookup deployments that point to a pod
func findDeploymentsForPod(pod *corev1.Pod, c cache.Cache) ([]*appsv1.Deployment, error) {
	var results []*appsv1.Deployment
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := cache.CacheKey{
		Namespace:  pod.Namespace,
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving deployments")
	}
	for _, u := range ul {
		d := &appsv1.Deployment{}
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, d)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured deployment")
		}
		if err := copyObjectMeta(d, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		labelSelector, err := getSelector(d)
		if err != nil {
			return nil, errors.Wrapf(err, "creating pod selector for deployment: %v", d.Name)
		}
		selector, err := metav1.LabelSelectorAsSelector(labelSelector)
		if err != nil {
			return nil, errors.Wrap(err, "invalid selector")
		}

		if selector.Empty() || !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}
		results = append(results, d)
	}
	return results, nil
}

func findDeploymentForReplicaSet(rs *appsv1.ReplicaSet, c cache.Cache) (*appsv1.Deployment, error) {
	if rs == nil {
		return nil, errors.New("nil replicaset")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key cache.CacheKey
	for _, owner := range rs.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "Deployment" {
			continue
		}

		key = cache.CacheKey{
			Namespace:  rs.Namespace,
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       owner.Name,
		}
	}
	if (key == cache.CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving deployment: %v", key)
	}
	for _, u := range ul {
		d := &appsv1.Deployment{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, d)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured deployment")
		}
		if err := copyObjectMeta(d, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		return d, nil
	}
	return nil, nil
}

func findStatefulSetForPod(pod *corev1.Pod, c cache.Cache) (*appsv1.StatefulSet, error) {
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key cache.CacheKey
	for _, owner := range pod.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "StatefulSet" {
			continue
		}

		key = cache.CacheKey{
			Namespace:  pod.Namespace,
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
			Name:       owner.Name,
		}
	}
	if (key == cache.CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving statefulset: %v", key)
	}
	for _, u := range ul {
		s := &appsv1.StatefulSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, s)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured statefulset")
		}
		if err := copyObjectMeta(s, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		return s, nil
	}
	return nil, nil
}

func findReplicationControllerForPod(pod *corev1.Pod, c cache.Cache) (*corev1.ReplicationController, error) {
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key cache.CacheKey
	for _, owner := range pod.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "ReplicationController" {
			continue
		}

		key = cache.CacheKey{
			Namespace:  pod.Namespace,
			APIVersion: "v1",
			Kind:       "ReplicationController",
			Name:       owner.Name,
		}
	}
	if (key == cache.CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving replicationcontroller: %v", key)
	}
	for _, u := range ul {
		rc := &corev1.ReplicationController{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, rc)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured replicationcontroller")
		}
		if err := copyObjectMeta(rc, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		return rc, nil
	}
	return nil, nil
}

func findDaemonSetForPod(pod *corev1.Pod, c cache.Cache) (*appsv1.DaemonSet, error) {
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key cache.CacheKey
	for _, owner := range pod.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "DaemonSet" {
			continue
		}

		key = cache.CacheKey{
			Namespace:  pod.Namespace,
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
			Name:       owner.Name,
		}
	}
	if (key == cache.CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving daemonset: %v", key)
	}
	for _, u := range ul {
		ds := &appsv1.DaemonSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, ds)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured daemonset")
		}
		if err := copyObjectMeta(ds, u); err != nil {
			return nil, errors.Wrap(err, "copying object metadata")
		}
		return ds, nil
	}
	return nil, nil
}

func listContains(lst []string, s string) bool {
	for _, item := range lst {
		if item == s {
			return true
		}
	}
	return false
}

func containsBackend(lst []v1beta1.IngressBackend, s string) bool {
	for _, item := range lst {
		if item.ServiceName == s {
			return true
		}
	}
	return false
}

type podGroup struct {
	UID    string
	Name   string
	Labels map[string]string
}

type podGroupKey struct {
	selector string
	ownerRef types.UID
}

func sortedLabels(labels map[string]string) string {
	keys := make([]string, len(labels))

	i := 0
	for k := range labels {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	// TODO hash instead?
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(labels[k])
		sb.WriteString(",")
	}
	return sb.String()
}

func podGroupKeyForPod(pod *corev1.Pod) podGroupKey {
	if pod == nil {
		return podGroupKey{}
	}

	var ownerRef types.UID
	if len(pod.OwnerReferences) > 0 {
		ownerRef = pod.OwnerReferences[0].UID
	} else {
		ownerRef = pod.UID
	}

	return podGroupKey{
		selector: sortedLabels(pod.Labels),
		ownerRef: ownerRef,
	}
}

func groupPods(pods []*corev1.Pod) []*podGroup {
	m := make(map[podGroupKey]bool)
	results := make([]*podGroup, 0, 1)
	for _, pod := range pods {
		k := podGroupKeyForPod(pod)
		if _, ok := m[k]; !ok {
			uid := fmt.Sprintf("pods-%s", string(k.ownerRef))
			grp := &podGroup{
				Name:   strings.TrimSuffix(k.selector, ","),
				UID:    uid,
				Labels: pod.Labels,
			}
			results = append(results, grp)
			m[k] = true
		}
	}
	return results
}

func serviceNames(backends []v1beta1.IngressBackend) []string {
	names := make([]string, 0, len(backends))
	for _, b := range backends {
		if b.ServiceName == "" {
			continue
		}
		names = append(names, b.ServiceName)
	}
	return names
}

func visitPod(ctx context.Context, pod *corev1.Pod, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if pod == nil {
		return errors.New("nil pod")
	}

	key := visitKeyForObject(pod)
	if visited[key] {
		return nil
	}
	visited[key] = true

	// Node is added by visitPodGroup, we will only explore the edges

	// Handle back-edges
	services, err := findServicesForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding services referencing pod: %v", pod.Name)
	}
	for _, svc := range services {
		if err = visitService(ctx, svc, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	replicaSets, err := findReplicaSetsForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding replicaSets referencing pod: %v", pod.Name)
	}
	for _, rs := range replicaSets {
		if err = visitReplicaSet(ctx, rs, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	deployments, err := findDeploymentsForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding deployments referencing pod: %v", pod.Name)
	}
	for _, d := range deployments {
		if err := visitDeployment(ctx, d, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	s, err := findStatefulSetForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding deployments referencing pod: %v", pod.Name)
	}
	if s != nil {
		if err = visitStatefulSet(ctx, s, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	rc, err := findReplicationControllerForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding replicationControllers referencing pod: %v", pod.Name)
	}
	if rc != nil {
		if err = visitReplicationController(ctx, rc, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	ds, err := findDaemonSetForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding daemonSet referencing pod: %v", pod.Name)
	}
	if ds != nil {
		if err = visitDaemonSet(ctx, ds, c, nodes, edges, visited); err != nil {
			return err
		}
	}
	return nil
}

func visitPodRoot(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil pod")
	}

	pod, ok := object.(*corev1.Pod)
	if !ok {
		return errors.Errorf("expected pod; received %T", object)
	}

	podList := &corev1.PodList{Items: []corev1.Pod{*pod}}
	return visitPodGroups(ctx, podList, nil, c, nodes, edges, visited)
}

func visitPodGroups(ctx context.Context, object runtime.Object, edgeFn edgeFunc, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil pod list")
	}

	podList, ok := object.(*corev1.PodList)
	if !ok {
		return errors.Errorf("expected pod list; received %T", object)
	}

	var pods []*corev1.Pod
	for _, pod := range podList.Items {
		pods = append(pods, &pod)
	}
	groups := groupPods(pods)

	for _, grp := range groups {
		if err := visitPodGroup(ctx, grp, c, nodes, edges, visited); err != nil {
			return err
		}

		if edgeFn != nil {
			edgeFn(string(grp.UID))
		}
	}

	for _, pod := range podList.Items {
		if err := visitPod(ctx, &pod, c, nodes, edges, visited); err != nil {
			return errors.Wrapf(err, "visiting pod %v", pod.Name)
		}
	}
	return nil
}

func visitPodGroup(ctx context.Context, grp *podGroup, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if grp == nil {
		return errors.New("nil podGroup")
	}

	key := visitKey{k8stypes.UID(grp.UID)}
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       grp.Name,
		APIVersion: "v1",
		Kind:       "Pods", // TODO podlist?
		Status:     statusForPodGroup(grp),
		IsNetwork:  false,
		Views:      []content.Content{},
	}
	uid := grp.UID
	nodes[uid] = node

	return nil
}

func visitService(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil service")
	}

	svc, ok := object.(*corev1.Service)
	if !ok {
		return errors.Errorf("expected service; received %T", object)
	}

	key := visitKeyForObject(svc)
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       svc.Name,
		APIVersion: svc.APIVersion,
		Kind:       svc.Kind,
		Status:     statusForService(svc),
		IsNetwork:  true,
		Views:      []content.Content{},
	}
	uid := string(svc.UID)
	nodes[uid] = node

	// Handle edges
	pods, err := findPodsForService(svc, c)
	if err != nil {
		return errors.Wrapf(err, "fetching pods for service %v", svc.Name)
	}

	podList := &corev1.PodList{}
	for _, pod := range pods {
		podList.Items = append(podList.Items, *pod)
	}

	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeImplicit, Node: dst})
	}
	if err = visitPodGroups(ctx, podList, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	// Reverse-lookup ingresses that reference the service
	ingresses, err := findIngressesForService(svc, c)
	if err != nil {
		return errors.Wrapf(err, "reverse-lookup ingresses for service %v", svc.Name)
	}
	for _, ingress := range ingresses {
		if err := visitIngress(ctx, ingress, c, nodes, edges, visited); err != nil {
			return errors.Wrapf(err, "visiting ingress for service %v", svc.Name)
		}
	}
	return nil
}

func visitReplicaSet(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil replica set")
	}

	rs, ok := object.(*appsv1.ReplicaSet)
	if !ok {
		return errors.Errorf("expected replica set; received %T", object)
	}

	key := visitKeyForObject(rs)
	if visited[key] {
		return nil
	}
	visited[key] = true

	checker := newNodeStatus(replicasSetCheckAvailableReplicas)
	statuses, err := checker.check(rs)
	if err != nil {
		return err
	}

	node := &content.Node{
		Name:       rs.Name,
		APIVersion: rs.APIVersion,
		Kind:       rs.Kind,
		Status:     statuses.Collapse(),
		IsNetwork:  false,
		Views:      []content.Content{},
	}
	uid := string(rs.UID)
	nodes[uid] = node

	// Handle edges
	pods, err := listPods(rs.GetNamespace(), rs.Spec.Selector, rs.UID, c)
	if err != nil {
		return errors.Wrapf(err, "fetching pods for replicaset %v", rs.Name)
	}
	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}

	podList := &corev1.PodList{}
	for _, pod := range pods {
		podList.Items = append(podList.Items, *pod)
	}

	if err = visitPodGroups(ctx, podList, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	// Handle back-edges
	d, err := findDeploymentForReplicaSet(rs, c)
	if err != nil {
		return errors.Wrapf(err, "finding deployment for replicaset %v", rs.Name)
	}
	if err := visitDeployment(ctx, d, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func visitIngress(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil ingress")
	}

	ingress, ok := object.(*v1beta1.Ingress)
	if !ok {
		return errors.Errorf("expected ingress; received %T", object)
	}

	key := visitKeyForObject(ingress)
	if visited[key] {
		return nil
	}
	visited[key] = true

	statusList, err := statusForIngress(ingress, c)
	if err != nil {
		return errors.Wrapf(err, "determining status for ingress %v", ingress.Name)
	}

	node := &content.Node{
		Name:       ingress.Name,
		APIVersion: ingress.APIVersion,
		Kind:       ingress.Kind,
		Status:     statusList.Collapse(),
		IsNetwork:  true,
		Views:      []content.Content{},
	}
	uid := string(ingress.UID)
	nodes[uid] = node

	// Handle edges
	backends, err := listIngressBackends(ingress, c)
	if err != nil {
		return errors.Wrapf(err, "listing backends for ingress %v", ingress.Name)
	}
	services, err := loadServices(serviceNames(backends), ingress.Namespace, c)
	if err != nil {
		return errors.Wrapf(err, "loading backends for ingress %v", ingress.Name)
	}
	for _, svc := range services {
		err := visitService(ctx, svc, c, nodes, edges, visited)
		if err != nil {
			return err
		}
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: string(svc.UID)})
	}
	return nil
}

func visitDeployment(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil deployment")
	}

	deployment, ok := object.(*appsv1.Deployment)
	if !ok {
		return errors.Errorf("expected deployment; received %T", object)
	}

	key := visitKeyForObject(deployment)
	if visited[key] {
		return nil
	}
	visited[key] = true

	checker := newNodeStatus(deploymentCheckUnavailable)
	statuses, err := checker.check(deployment)
	if err != nil {
		return err
	}

	node := &content.Node{
		Name:       deployment.Name,
		APIVersion: deployment.APIVersion,
		Kind:       deployment.Kind,
		Status:     statuses.Collapse(),
		IsNetwork:  false,
		Views:      []content.Content{},
	}
	uid := string(deployment.UID)
	nodes[uid] = node

	// Handle edges
	rsList, err := listReplicaSets(deployment, c)
	if err != nil {
		return errors.Wrapf(err, "fetching replicasets for deployment %v", deployment.Name)
	}

	var currentReplicaSets []*appsv1.ReplicaSet

	if rs := findNewReplicaSet(deployment, rsList); rs != nil {
		currentReplicaSets = append(currentReplicaSets, rs)
	}

	for _, rs := range findOldReplicaSets(deployment, rsList) {
		currentReplicaSets = append(currentReplicaSets, rs)
	}
	for _, rs := range currentReplicaSets {

		err := visitReplicaSet(ctx, rs, c, nodes, edges, visited)
		if err != nil {
			return err
		}
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: string(rs.UID)})
	}

	return nil
}

func visitStatefulSet(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil stateful set")
	}

	s, ok := object.(*appsv1.StatefulSet)
	if !ok {
		return errors.Errorf("expected stateful set; received %T", object)
	}

	key := visitKeyForObject(s)
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       s.Name,
		APIVersion: s.APIVersion,
		Kind:       s.Kind,
		Status:     statusForStatefulSet(s),
		IsNetwork:  false,
		Views:      []content.Content{},
	}
	uid := string(s.UID)
	nodes[uid] = node

	// Handle edges
	pods, err := listPods(s.GetNamespace(), s.Spec.Selector, s.UID, c)
	if err != nil {
		return errors.Wrapf(err, "fetching pods for statefulset %v", s.Name)
	}

	podList := &corev1.PodList{}
	for _, pod := range pods {
		podList.Items = append(podList.Items, *pod)
	}

	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}
	if err := visitPodGroups(ctx, podList, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func visitReplicationController(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil replicationcontroller")
	}

	rc, ok := object.(*corev1.ReplicationController)
	if !ok {
		return errors.Errorf("expected replicationcontroller; received %T", object)
	}

	key := visitKeyForObject(rc)
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       rc.Name,
		APIVersion: rc.APIVersion,
		Kind:       rc.Kind,
		Status:     statusForReplicationController(rc),
		IsNetwork:  false,
		Views:      []content.Content{},
	}
	uid := string(rc.UID)
	nodes[uid] = node

	// Handle edges
	selector, err := getSelector(rc)
	if err != nil {
		return errors.Wrapf(err, "fetching selector for replicationcontroller: %v", rc.Name)
	}
	pods, err := listPods(rc.GetNamespace(), selector, rc.UID, c)
	if err != nil {
		return errors.Wrapf(err, "fetching pods for replicationcontroller %v", rc.Name)
	}

	podList := &corev1.PodList{}
	for _, pod := range pods {
		podList.Items = append(podList.Items, *pod)
	}

	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}
	if err := visitPodGroups(ctx, podList, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func visitDaemonSet(ctx context.Context, object runtime.Object, c cache.Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if object == nil {
		return errors.New("nil daemonset")
	}

	ds, ok := object.(*appsv1.DaemonSet)
	if !ok {
		return errors.Errorf("expected daemonset; received %T", object)
	}

	key := visitKeyForObject(ds)
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       ds.Name,
		APIVersion: ds.APIVersion,
		Kind:       ds.Kind,
		Status:     statusForDaemonSet(ds),
		IsNetwork:  false,
		Views:      []content.Content{},
	}
	uid := string(ds.UID)
	nodes[uid] = node

	// Handle edges
	pods, err := listPods(ds.GetNamespace(), ds.Spec.Selector, ds.UID, c)
	if err != nil {
		return errors.Wrapf(err, "fetching pods for daemonset %v", ds.Name)
	}

	podList := &corev1.PodList{}
	for _, pod := range pods {
		podList.Items = append(podList.Items, *pod)
	}

	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}
	if err := visitPodGroups(ctx, podList, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}
