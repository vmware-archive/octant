package printer

import (
	"fmt"
	"path"
	"strings"

	"github.com/heptio/developer-dash/internal/portforward"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// ContainerConfiguration generates container configuration.
type ContainerConfiguration struct {
	parent             runtime.Object
	container          *corev1.Container
	portForwardService portforward.PortForwardInterface
	isInit             bool
}

// NewContainerConfiguration creates an instance of ContainerConfiguration.
func NewContainerConfiguration(parent runtime.Object, c *corev1.Container, pfs portforward.PortForwardInterface, isInit bool) *ContainerConfiguration {
	return &ContainerConfiguration{
		parent:             parent,
		container:          c,
		isInit:             isInit,
		portForwardService: pfs,
	}
}

// Create creates a deployment configuration summary.
func (cc *ContainerConfiguration) Create() (*component.Summary, error) {
	if cc == nil || cc.container == nil {
		return nil, errors.New("container is nil")
	}
	c := cc.container

	sections := component.SummarySections{}

	sections.AddText("Image", c.Image)

	hostPorts := describeContainerHostPorts(c.Ports)
	if hostPorts != "" {
		sections.AddText("Host Ports", hostPorts)
	}
	containerPorts, err := describeContainerPorts(cc.parent, c.Ports, cc.portForwardService)
	if err != nil {
		return nil, errors.Wrap(err, "describe container ports")
	}
	if len(containerPorts) > 0 {
		sections.Add("Container Ports", component.NewPorts(containerPorts))
	}

	envTbl, err := describeContainerEnv(cc.parent, c)
	if err != nil {
		return nil, errors.Wrap(err, "describing environment")
	}

	if len(envTbl.Config.Rows) > 0 {
		sections.Add("Environment", envTbl)
	}

	cmd := printSlice(c.Command)
	if cmd != "" {
		sections.AddText("Command", cmd)
	}
	args := printSlice(c.Args)
	if args != "" {
		sections.AddText("Args", args)
	}

	if len(c.VolumeMounts) > 0 {
		sections.Add("Volume Mounts", describeVolumeMounts(c))
	}

	title := "Container"
	if cc.isInit {
		title = "Init Container"
	}

	summary := component.NewSummary(fmt.Sprintf("%s %s", title, c.Name), sections...)
	return summary, nil
}

func isPodGVK(gvk schema.GroupVersionKind) bool {
	return gvk.Group == "" && gvk.Version == "v1" && gvk.Kind == "Pod"
}

type notFound interface {
	NotFound() bool
}

func isNotFound(err error) bool {
	if err == nil {
		return false
	}

	notFoundErr, ok := errors.Cause(err).(notFound)
	if !ok {
		return false
	}

	return notFoundErr.NotFound()
}

func describeContainerPorts(
	parent runtime.Object,
	cPorts []corev1.ContainerPort,
	portForwardService portforward.PortForwardInterface) ([]component.Port, error) {
	var list []component.Port

	var namespace string
	var name string
	var err error
	gvk := parent.GetObjectKind().GroupVersionKind()
	isPod := isPodGVK(gvk)
	if isPod {
		accessor := meta.NewAccessor()
		namespace, err = accessor.Namespace(parent)
		if err != nil {
			return nil, errors.Wrap(err, "find parent namespace")
		}

		name, err = accessor.Name(parent)
		if err != nil {
			return nil, errors.Wrap(err, "find parent name")
		}
	}

	for _, cPort := range cPorts {
		if cPort.ContainerPort == 0 {
			continue
		}
		pfs := component.PortForwardState{}

		var port *component.Port
		if isPod && cPort.Protocol == "TCP" {
			pfs.IsForwardable = true
			state, err := portForwardService.Find(namespace, gvk, name)
			if err != nil {
				if _, ok := err.(notFound); !ok {
					return nil, errors.Wrap(err, "query portforward service for pod")
				}
			} else {
				pfs.ID = state.ID
				for _, forwarded := range state.Ports {
					if int(forwarded.Remote) == int(cPort.ContainerPort) {
						pfs.Port = int(forwarded.Local)
					}
				}
				pfs.IsForwarded = true
			}
		}

		apiVersion, kind := gvk.ToAPIVersionAndKind()

		port = component.NewPort(
			namespace,
			apiVersion,
			kind,
			name,
			int(cPort.ContainerPort),
			string(cPort.Protocol), pfs)
		list = append(list, *port)
	}
	return list, nil
}

func describeContainerHostPorts(cPorts []corev1.ContainerPort) string {
	ports := make([]string, 0, len(cPorts))
	for _, cPort := range cPorts {
		if cPort.HostPort == 0 {
			continue
		}
		ports = append(ports, fmt.Sprintf("%d/%s", cPort.HostPort, cPort.Protocol))
	}
	return strings.Join(ports, ", ")
}

// describeContainerEnv returns a table describing a container environment
func describeContainerEnv(parent runtime.Object, c *corev1.Container) (*component.Table, error) {
	if c == nil {
		return nil, errors.New("container is nil")
	}

	var ns string
	if parent != nil {
		accessor := meta.NewAccessor()
		ns, _ = accessor.Namespace(parent)
	}
	if ns == "" {
		ns = "default"
	}

	cols := component.NewTableCols("Name", "Value", "Source")
	tbl := component.NewTable("Environment", cols)
	tbl.Add(describeEnvRows(ns, c.Env)...)
	tbl.Add(describeEnvFromRows(ns, c.EnvFrom)...)
	return tbl, nil
}

// describeEnvRows renders container environment variables as table rows.
// Expected columns: Name, Value, Source
func describeEnvRows(namespace string, vars []corev1.EnvVar) []component.TableRow {
	rows := make([]component.TableRow, 0)
	for _, e := range vars {
		row := component.TableRow{}
		rows = append(rows, row)

		// each row requires data for each defined column
		row["Source"] = component.NewText("")

		row["Name"] = component.NewText(e.Name)
		row["Value"] = component.NewText(e.Value) // TODO should we resolve values when source is a reference (valueFrom)?
		if e.Value != "" || e.ValueFrom == nil {
			continue
		}

		switch {
		case e.ValueFrom.FieldRef != nil:
			ref := e.ValueFrom.FieldRef
			row["Source"] = component.NewText(ref.FieldPath)
		case e.ValueFrom.ResourceFieldRef != nil:
			ref := e.ValueFrom.ResourceFieldRef
			row["Source"] = component.NewText(ref.Resource)
		case e.ValueFrom.SecretKeyRef != nil:
			ref := e.ValueFrom.SecretKeyRef
			row["Source"] = link.ForGVK(namespace, "v1", "Secret", ref.Name,
				fmt.Sprintf("%s:%s", ref.Name, ref.Key))
		case e.ValueFrom.ConfigMapKeyRef != nil:
			ref := e.ValueFrom.ConfigMapKeyRef
			row["Source"] = link.ForGVK(namespace, "v1", "ConfigMap", ref.Name,
				fmt.Sprintf("%s:%s", ref.Name, ref.Key))
		}
	}

	return rows
}

// describeEnvFromRows renders container environmentFrom references as table rows.
// Expected columns: Name, Value, Source
// TODO: Consider expanding variables from referenced configmap / secret
func describeEnvFromRows(namespace string, vars []corev1.EnvFromSource) []component.TableRow {
	rows := make([]component.TableRow, 0)
	for _, e := range vars {
		row := component.TableRow{}
		rows = append(rows, row)

		switch {
		case e.SecretRef != nil:
			ref := e.SecretRef
			row["Source"] = link.ForGVK(namespace, "v1", "Secret", ref.Name, ref.Name)
		case e.ConfigMapRef != nil:
			ref := e.ConfigMapRef
			row["Source"] = link.ForGVK(namespace, "v1", "ConfigMap", ref.Name, ref.Name)
		}
	}

	return rows
}

// printSlice returns a string representation of a string slice, in a format similar
// to ['a', 'b', 'c']. An empty slice will be returned as an empty string, rather than [].
func printSlice(s []string) string {
	if len(s) == 0 {
		return ""
	}

	return fmt.Sprintf("['%s']", strings.Join(s, "', '"))
}

func printVolumeMountPath(mnt corev1.VolumeMount) string {
	var access string
	if mnt.ReadOnly {
		access = "ro"
	} else {
		access = "rw"
	}
	p := path.Join(mnt.MountPath, mnt.SubPath)

	return fmt.Sprintf("%s (%s)", p, access)
}

func describeVolumeMounts(c *corev1.Container) *component.Table {
	cols := component.NewTableCols("Name", "Mount Path", "Propagation")
	tbl := component.NewTable("Volume Mounts", cols)

	for _, m := range c.VolumeMounts {
		row := component.TableRow{}
		row["Name"] = component.NewText(m.Name)
		row["Mount Path"] = component.NewText(printVolumeMountPath(m))

		var prop string
		if m.MountPropagation != nil {
			prop = string(*m.MountPropagation)
		}
		row["Propagation"] = component.NewText(prop)

		tbl.Add(row)
	}

	return tbl
}
