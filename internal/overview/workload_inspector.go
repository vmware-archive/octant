package overview

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

// WorkloadInspector is both a View and a View
type workloadInspectorView struct {
}

type visitKey struct {
	uid k8stypes.UID
}
type visitSet map[visitKey]bool

// Lets us integrate as a table in a Resource ObjectView
func newWorkloadInspectorView(prefix, namespace string, c clock.Clock) View {
	return &workloadInspectorView{}
}

// Implements View.Content
func (wid *workloadInspectorView) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	visited := visitSet{}
	dag := content.NewDAG()

	acc := meta.NewAccessor()
	uid, err := acc.UID(object)
	if err != nil {
		return nil, errors.Wrapf(err, "fetching uid for object, type: %T", object)
	}
	dag.Selected = string(uid)

	switch v := object.(type) {
	case (*core.Pod):
		if err := wid.visitPodGroups(ctx, []*core.Pod{v}, nil, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*core.Service):
		if err := wid.visitService(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*extensions.Deployment):
		if err := wid.visitDeployment(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*extensions.ReplicaSet):
		if err := wid.visitReplicaSet(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*v1beta1.Ingress):
		if err := wid.visitIngress(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*apps.StatefulSet):
		if err := wid.visitStatefulSet(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*core.ReplicationController):
		if err := wid.visitReplicationController(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	case (*extensions.DaemonSet):
		if err := wid.visitDaemonSet(ctx, v, c, dag.Nodes, dag.Edges, visited); err != nil {
			return nil, err
		}
	default:
	}

	return []content.Content{dag}, nil
}

func visitKeyForObject(obj runtime.Object) visitKey {
	acc := meta.NewAccessor()
	// TODO ERROR HANDLING
	uid, _ := acc.UID(obj)
	return visitKey{uid}
}

func (wid *workloadInspectorView) visitPod(ctx context.Context, pod *core.Pod, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
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
		if err := wid.visitService(ctx, svc, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	replicaSets, err := findReplicaSetsForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding replicaSets referencing pod: %v", pod.Name)
	}
	for _, rs := range replicaSets {
		if err := wid.visitReplicaSet(ctx, rs, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	deployments, err := findDeploymentsForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding deployments referencing pod: %v", pod.Name)
	}
	for _, d := range deployments {
		if err := wid.visitDeployment(ctx, d, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	s, err := findStatefulSetForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding deployments referencing pod: %v", pod.Name)
	}
	if s != nil {
		if err := wid.visitStatefulSet(ctx, s, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	rc, err := findReplicationControllerForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding replicationControllers referencing pod: %v", pod.Name)
	}
	if rc != nil {
		if err := wid.visitReplicationController(ctx, rc, c, nodes, edges, visited); err != nil {
			return err
		}
	}

	ds, err := findDaemonSetForPod(pod, c)
	if err != nil {
		return errors.Wrapf(err, "finding daemonSet referencing pod: %v", pod.Name)
	}
	if ds != nil {
		if err := wid.visitDaemonSet(ctx, ds, c, nodes, edges, visited); err != nil {
			return err
		}
	}
	return nil
}

// An edgeFunc will create an edge to the provided destination node
type edgeFunc func(dst string)

func (wid *workloadInspectorView) visitPodGroups(ctx context.Context, pods []*core.Pod, edgeFn edgeFunc, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	groups := groupPods(pods)

	for _, grp := range groups {
		wid.visitPodGroup(ctx, grp, c, nodes, edges, visited)

		if edgeFn != nil {
			edgeFn(string(grp.UID))
		}
	}

	for _, pod := range pods {
		if err := wid.visitPod(ctx, pod, c, nodes, edges, visited); err != nil {
			return errors.Wrapf(err, "visiting pod %v", pod.Name)
		}
	}
	return nil
}

func (wid *workloadInspectorView) visitPodGroup(ctx context.Context, grp *podGroup, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
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

func (wid *workloadInspectorView) visitDeployment(ctx context.Context, deployment *extensions.Deployment, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if deployment == nil {
		return errors.New("nil deployment")
	}

	key := visitKeyForObject(deployment)
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       deployment.Name,
		APIVersion: deployment.APIVersion,
		Kind:       deployment.Kind,
		Status:     statusForDeployment(deployment),
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

	var currentReplicaSets []*extensions.ReplicaSet

	if rs := findNewReplicaSet(deployment, rsList); rs != nil {
		currentReplicaSets = append(currentReplicaSets, rs)
	}

	for _, rs := range findOldReplicaSets(deployment, rsList) {
		currentReplicaSets = append(currentReplicaSets, rs)
	}

	for _, rs := range currentReplicaSets {
		err := wid.visitReplicaSet(ctx, rs, c, nodes, edges, visited)
		if err != nil {
			return err
		}
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: string(rs.UID)})
	}

	return nil
}

func (wid *workloadInspectorView) visitReplicaSet(ctx context.Context, rs *extensions.ReplicaSet, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if rs == nil {
		return errors.New("nil replicaset")
	}

	key := visitKeyForObject(rs)
	if visited[key] {
		return nil
	}
	visited[key] = true

	node := &content.Node{
		Name:       rs.Name,
		APIVersion: rs.APIVersion,
		Kind:       rs.Kind,
		Status:     statusForReplicaSet(rs),
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
	if err := wid.visitPodGroups(ctx, pods, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	// Handle back-edges
	d, err := findDeploymentForReplicaSet(rs, c)
	if err != nil {
		return errors.Wrapf(err, "finding deployment for replicaset %v", rs.Name)
	}
	if err := wid.visitDeployment(ctx, d, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func (wid *workloadInspectorView) visitService(ctx context.Context, svc *core.Service, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if svc == nil {
		return errors.New("nil service")
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
	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeImplicit, Node: dst})
	}
	if err := wid.visitPodGroups(ctx, pods, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	// Reverse-lookup ingresses that reference the service
	ingresses, err := findIngressesForService(svc, c)
	if err != nil {
		return errors.Wrapf(err, "reverse-lookup ingresses for service %v", svc.Name)
	}
	for _, ingress := range ingresses {
		if err := wid.visitIngress(ctx, ingress, c, nodes, edges, visited); err != nil {
			return errors.Wrapf(err, "visiting ingress for service %v", svc.Name)
		}
	}
	return nil
}

func (wid *workloadInspectorView) visitIngress(ctx context.Context, ingress *v1beta1.Ingress, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if ingress == nil {
		return errors.New("nil ingress")
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
		err := wid.visitService(ctx, svc, c, nodes, edges, visited)
		if err != nil {
			return err
		}
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: string(svc.UID)})
	}
	return nil
}

func (wid *workloadInspectorView) visitStatefulSet(ctx context.Context, s *apps.StatefulSet, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if s == nil {
		return errors.New("nil statefulset")
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
	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}
	if err := wid.visitPodGroups(ctx, pods, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func (wid *workloadInspectorView) visitReplicationController(ctx context.Context, rc *core.ReplicationController, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if rc == nil {
		return errors.New("nil replicationcontroller")
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
	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}
	if err := wid.visitPodGroups(ctx, pods, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func (wid *workloadInspectorView) visitDaemonSet(ctx context.Context, ds *extensions.DaemonSet, c Cache, nodes content.Nodes, edges content.AdjList, visited visitSet) error {
	if ds == nil {
		return errors.New("nil daemonset")
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
	newEdgeFn := func(dst string) {
		edges.Add(uid, content.Edge{Type: content.EdgeTypeExplicit, Node: dst})
	}
	if err := wid.visitPodGroups(ctx, pods, newEdgeFn, c, nodes, edges, visited); err != nil {
		return err
	}

	return nil
}

func listIngressPaths(ingress *v1beta1.Ingress, c Cache) ([]v1beta1.HTTPIngressPath, error) {
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
func listIngressBackends(ingress *v1beta1.Ingress, c Cache) ([]v1beta1.IngressBackend, error) {
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

func loadServices(serviceNames []string, namespace string, c Cache) ([]*core.Service, error) {
	var services []*core.Service
	for _, backend := range serviceNames {
		key := CacheKey{
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
			svc := &core.Service{}
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

func loadService(name string, namespace string, c Cache) (*core.Service, error) {
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
func findIngressesForService(svc *core.Service, c Cache) ([]*v1beta1.Ingress, error) {
	var results []*v1beta1.Ingress
	if svc == nil {
		return nil, errors.New("nil service")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := CacheKey{
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
		if err := copyObjectMeta(ingress, u); err != nil {
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

func findPodsForService(svc *core.Service, c Cache) ([]*core.Pod, error) {
	if svc == nil {
		return nil, errors.New("nil service")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}
	key := CacheKey{
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
func findServicesForPod(pod *core.Pod, c Cache) ([]*core.Service, error) {
	var results []*core.Service
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := CacheKey{
		Namespace:  pod.Namespace,
		APIVersion: "v1",
		Kind:       "Service",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving services")
	}
	for _, u := range ul {
		svc := &core.Service{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, svc)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured service")
		}
		if err := copyObjectMeta(svc, u); err != nil {
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
func findReplicaSetsForPod(pod *core.Pod, c Cache) ([]*extensions.ReplicaSet, error) {
	var results []*extensions.ReplicaSet
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := CacheKey{
		Namespace:  pod.Namespace,
		APIVersion: "apps/v1",
		Kind:       "ReplicaSet",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving replicaSets")
	}
	for _, u := range ul {
		rs := &extensions.ReplicaSet{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, rs)
		if err != nil {
			return nil, errors.Wrap(err, "converting unstructured replicaSet")
		}
		if err := copyObjectMeta(rs, u); err != nil {
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
func findDeploymentsForPod(pod *core.Pod, c Cache) ([]*extensions.Deployment, error) {
	var results []*extensions.Deployment
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	key := CacheKey{
		Namespace:  pod.Namespace,
		APIVersion: "apps/v1",
		Kind:       "Deployment",
	}
	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving deployments")
	}
	for _, u := range ul {
		d := &extensions.Deployment{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, d)
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

func findDeploymentForReplicaSet(rs *extensions.ReplicaSet, c Cache) (*extensions.Deployment, error) {
	if rs == nil {
		return nil, errors.New("nil replicaset")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key CacheKey
	for _, owner := range rs.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "Deployment" {
			continue
		}

		key = CacheKey{
			Namespace:  rs.Namespace,
			APIVersion: "apps/v1",
			Kind:       "Deployment",
			Name:       owner.Name,
		}
	}
	if (key == CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving deployment: %v", key)
	}
	for _, u := range ul {
		d := &extensions.Deployment{}
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

func findStatefulSetForPod(pod *core.Pod, c Cache) (*apps.StatefulSet, error) {
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key CacheKey
	for _, owner := range pod.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "StatefulSet" {
			continue
		}

		key = CacheKey{
			Namespace:  pod.Namespace,
			APIVersion: "apps/v1",
			Kind:       "StatefulSet",
			Name:       owner.Name,
		}
	}
	if (key == CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving statefulset: %v", key)
	}
	for _, u := range ul {
		s := &apps.StatefulSet{}
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

func findReplicationControllerForPod(pod *core.Pod, c Cache) (*core.ReplicationController, error) {
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key CacheKey
	for _, owner := range pod.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "ReplicationController" {
			continue
		}

		key = CacheKey{
			Namespace:  pod.Namespace,
			APIVersion: "v1",
			Kind:       "ReplicationController",
			Name:       owner.Name,
		}
	}
	if (key == CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving replicationcontroller: %v", key)
	}
	for _, u := range ul {
		rc := &core.ReplicationController{}
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

func findDaemonSetForPod(pod *core.Pod, c Cache) (*extensions.DaemonSet, error) {
	if pod == nil {
		return nil, errors.New("nil pod")
	}
	if c == nil {
		return nil, errors.New("nil cache")
	}

	var key CacheKey
	for _, owner := range pod.OwnerReferences {
		// Don't compare APIVersion - there may be several aliases
		if owner.Kind != "DaemonSet" {
			continue
		}

		key = CacheKey{
			Namespace:  pod.Namespace,
			APIVersion: "apps/v1",
			Kind:       "DaemonSet",
			Name:       owner.Name,
		}
	}
	if (key == CacheKey{}) {
		return nil, nil
	}

	ul, err := c.Retrieve(key)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving daemonset: %v", key)
	}
	for _, u := range ul {
		ds := &extensions.DaemonSet{}
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
	// selector metav1.LabelSelector
	// selector map[string]string
	selector string
	// metav1.LabelSelectorAsMap
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

func podGroupKeyForPod(pod *core.Pod) podGroupKey {
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

func groupPods(pods []*core.Pod) []*podGroup {
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
