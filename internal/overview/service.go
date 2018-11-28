package overview

import (
	"context"
	"fmt"

	"github.com/heptio/developer-dash/internal/content"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/scheme"
)

type ServiceSummary struct{}

var _ View = (*ServiceSummary)(nil)

func NewServiceSummary(prefix, namespace string, c clock.Clock) View {
	return &ServiceSummary{}
}

func (js *ServiceSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	ss, err := retrieveService(object)
	if err != nil {
		return nil, err
	}

	detail, err := printServiceSummary(ss)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type ServicePort struct{}

var _ View = (*ServicePort)(nil)

func NewServicePort(prefix, namespace string, c clock.Clock) View {
	return &ServicePort{}
}

func (js *ServicePort) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	ss, err := retrieveService(object)
	if err != nil {
		return nil, err
	}

	portList := content.NewSummary("Ports", []content.Section{})

	for _, port := range ss.Spec.Ports {
		name := port.Name
		if name == "" {
			name = "<unnamed>"
		}

		section := content.NewSection()
		section.Title = name

		section.AddText("Port", fmt.Sprintf("%d/%s", port.Port, port.Protocol))
		if port.TargetPort.Type == intstr.Type(intstr.Int) {
			section.AddText("TargetPort", fmt.Sprintf("%d/%s", port.TargetPort.IntVal, port.Protocol))
		} else {
			section.AddText("TargetPort", fmt.Sprintf("%s/%s", port.TargetPort.StrVal, port.Protocol))
		}
		if port.NodePort != 0 {
			section.AddText("NodePort", fmt.Sprintf("%d/%s", port.NodePort, port.Protocol))
		}

		portList.Sections = append(portList.Sections, section)
	}

	return []content.Content{
		&portList,
	}, nil
}

var serviceEndpointsColumns = []string{
	"Host",
	"Ports (Name/Port/Protocol)",
	"Node",
	"Ready",
}

type ServiceEndpoints struct{}

var _ View = (*ServiceEndpoints)(nil)

func NewServiceEndpoints(prefix, namespace string, c clock.Clock) View {
	return &ServiceEndpoints{}
}

func (js *ServiceEndpoints) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	ss, err := retrieveService(object)
	if err != nil {
		return nil, err
	}

	endpoints, err := listEndpoints(ss.GetNamespace(), ss.GetName(), c)
	if err != nil {
		return nil, err
	}

	table := content.NewTable("Endpoints", "This service does not have any endpoints")
	for _, name := range serviceEndpointsColumns {
		table.Columns = append(table.Columns, tableCol(name))
	}

	podKey := CacheKey{
		Namespace:  ss.GetNamespace(),
		APIVersion: "v1",
		Kind:       "Pod",
	}

	pods, err := loadPods(podKey, c, nil)
	if err != nil {
		return nil, err
	}

	podMap := make(map[string]*core.Pod)

	for _, pod := range pods {
		podMap[pod.Spec.Hostname] = pod
	}

	for _, endpoint := range endpoints {
		for _, subset := range endpoint.Subsets {
			for _, port := range subset.Ports {
				portStr := fmt.Sprintf("%s / %d / %s",
					port.Name, port.Port, port.Protocol,
				)

				for _, address := range subset.Addresses {
					table.AddRow(endpointAddressRow(address, portStr, true))
				}

				for _, address := range subset.NotReadyAddresses {
					table.AddRow(endpointAddressRow(address, portStr, false))
				}

			}

		}
	}

	return []content.Content{
		&table,
	}, nil
}

func endpointAddressRow(address core.EndpointAddress, port string, ready bool) content.TableRow {
	nodeName := "<unset>"
	if address.NodeName != nil {
		nodeName = *address.NodeName
	}

	return content.TableRow{
		serviceEndpointsColumns[0]: content.NewStringText(address.IP),
		serviceEndpointsColumns[1]: content.NewStringText(port),
		serviceEndpointsColumns[2]: content.NewStringText(nodeName),
		serviceEndpointsColumns[3]: content.NewStringText(fmt.Sprintf("%t", ready)),
	}

}

func retrieveService(object runtime.Object) (*core.Service, error) {
	rc, ok := object.(*core.Service)
	if !ok {
		return nil, errors.Errorf("expected object to be a Service, it was %T", object)
	}

	return rc, nil
}

func listEndpoints(namespace string, name string, c Cache) ([]*core.Endpoints, error) {
	key := CacheKey{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Endpoints",
		Name:       name,
	}

	return loadEndpoints(key, c)
}

func loadEndpoints(key CacheKey, c Cache) ([]*core.Endpoints, error) {
	objects, err := c.Retrieve(key)
	if err != nil {
		return nil, err
	}

	var list []*core.Endpoints

	for _, object := range objects {
		e := &core.Endpoints{}
		if err := scheme.Scheme.Convert(object, e, runtime.InternalGroupVersioner); err != nil {
			return nil, err
		}

		if err := copyObjectMeta(e, object); err != nil {
			return nil, err
		}

		list = append(list, e)
	}

	return list, nil
}
