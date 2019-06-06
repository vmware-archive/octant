package printer

import (
	"fmt"
	"path"
	"strings"

	"github.com/heptio/developer-dash/internal/portforward"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/developer-dash/pkg/view/component"
)

// ContainerConfiguration generates container configuration.
type ContainerConfiguration struct {
	parent             runtime.Object
	container          *corev1.Container
	portForwardService portforward.PortForwarder
	isInit             bool
	options            Options
}

// NewContainerConfiguration creates an instance of ContainerConfiguration.
func NewContainerConfiguration(parent runtime.Object, c *corev1.Container, pfs portforward.PortForwarder, isInit bool, options Options) *ContainerConfiguration {
	return &ContainerConfiguration{
		parent:             parent,
		container:          c,
		isInit:             isInit,
		portForwardService: pfs,
		options:            options,
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

	if pod, ok := cc.parent.(*corev1.Pod); ok {
		status, err := findContainerStatus(pod, cc.container.Name, cc.isInit)
		if err != nil {
			return nil, errors.Wrapf(err, "get container status for %q", cc.container.Name)
		}

		sections.AddText("Last State", printContainerState(status.LastTerminationState))
		sections.AddText("Current State", printContainerState(status.State))
		sections.AddText("Ready", fmt.Sprintf("%t", status.Ready))
		sections.AddText("Restart Count", fmt.Sprintf("%d", status.RestartCount))
	}

	envTbl, err := describeContainerEnv(cc.parent, c, cc.options)
	if err != nil {
		return nil, errors.Wrap(err, "describing environment")
	}

	if len(envTbl.Rows()) > 0 {
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

func printContainerState(state corev1.ContainerState) string {
	switch {
	case state.Running != nil:
		return fmt.Sprintf("started at %s", state.Running.StartedAt)
	case state.Waiting != nil:
		return fmt.Sprintf("waiting: %s", state.Waiting.Message)
	case state.Terminated != nil:
		return fmt.Sprintf("terminated with %d at %s: %s",
			state.Terminated.ExitCode,
			state.Terminated.FinishedAt,
			state.Terminated.Reason)
	}

	return "indeterminate"
}

type containerNotFoundError struct {
	name string
}

func (e *containerNotFoundError) Error() string {
	return fmt.Sprintf("container %q not found", e.name)
}

func findContainerStatus(pod *corev1.Pod, name string, isInit bool) (*corev1.ContainerStatus, error) {
	if pod == nil {
		return nil, errors.New("pod is nil")
	}

	statuses := pod.Status.ContainerStatuses
	if isInit {
		statuses = pod.Status.InitContainerStatuses
	}

	for _, status := range statuses {
		if status.Name == name {
			return &status, nil
		}
	}

	return nil, &containerNotFoundError{name: name}
}

func isPodGVK(gvk schema.GroupVersionKind) bool {
	return gvk.Group == "" && gvk.Version == "v1" && gvk.Kind == "Pod"
}

type notFound interface {
	NotFound() bool
}

func describeContainerPorts(
	parent runtime.Object,
	cPorts []corev1.ContainerPort,
	portForwardService portforward.PortForwarder) ([]component.Port, error) {
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
					return nil, errors.Wrap(err, "query port forward service for pod")
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
func describeContainerEnv(parent runtime.Object, c *corev1.Container, options Options) (*component.Table, error) {
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

	envRows, err := describeEnvRows(ns, c.Env, options)
	if err != nil {
		return nil, err
	}
	tbl.Add(envRows...)
	envFromRows, err := describeEnvFromRows(ns, c.EnvFrom, options)
	if err != nil {
		return nil, err
	}
	tbl.Add(envFromRows...)
	return tbl, nil
}

// describeEnvRows renders container environment variables as table rows.
// Expected columns: Name, Value, Source
func describeEnvRows(namespace string, vars []corev1.EnvVar, options Options) ([]component.TableRow, error) {
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
			source, err := options.Link.ForGVK(namespace, "v1", "Secret", ref.Name,
				fmt.Sprintf("%s:%s", ref.Name, ref.Key))
			if err != nil {
				return nil, err
			}
			row["Source"] = source
		case e.ValueFrom.ConfigMapKeyRef != nil:
			ref := e.ValueFrom.ConfigMapKeyRef
			source, err := options.Link.ForGVK(namespace, "v1", "ConfigMap", ref.Name,
				fmt.Sprintf("%s:%s", ref.Name, ref.Key))
			if err != nil {
				return nil, err
			}
			row["Source"] = source
		}
	}

	return rows, nil
}

// describeEnvFromRows renders container environmentFrom references as table rows.
// Expected columns: Name, Value, Source
// TODO: Consider expanding variables from referenced config map / secret
func describeEnvFromRows(namespace string, vars []corev1.EnvFromSource, options Options) ([]component.TableRow, error) {
	rows := make([]component.TableRow, 0)
	for _, e := range vars {
		row := component.TableRow{}
		rows = append(rows, row)

		switch {
		case e.SecretRef != nil:
			ref := e.SecretRef
			source, err := options.Link.ForGVK(namespace, "v1", "Secret", ref.Name, ref.Name)
			if err != nil {
				return nil, err
			}
			row["Source"] = source
		case e.ConfigMapRef != nil:
			ref := e.ConfigMapRef
			source, err := options.Link.ForGVK(namespace, "v1", "ConfigMap", ref.Name, ref.Name)
			if err != nil {
				return nil, err
			}
			row["Source"] = source
		}
	}

	return rows, nil
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
