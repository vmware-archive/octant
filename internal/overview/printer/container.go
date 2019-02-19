package printer

import (
	"fmt"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

// ContainerConfiguration generates container configuration.
type ContainerConfiguration struct {
	parent    runtime.Object
	container *corev1.Container
}

// NewContainerConfiguration creates an instance of ContainerConfiguration.
func NewContainerConfiguration(parent runtime.Object, c *corev1.Container) *ContainerConfiguration {
	return &ContainerConfiguration{
		parent:    parent,
		container: c,
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
	containerPorts := describeContainerPorts(c.Ports)
	if containerPorts != "" {
		sections.AddText("Container Ports", containerPorts)
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

	summary := component.NewSummary(fmt.Sprintf("Container %s", c.Name), sections...)
	return summary, nil
}

func describeContainerPorts(cPorts []corev1.ContainerPort) string {
	ports := make([]string, 0, len(cPorts))
	for _, cPort := range cPorts {
		if cPort.ContainerPort == 0 {
			continue
		}
		ports = append(ports, fmt.Sprintf("%d/%s", cPort.ContainerPort, cPort.Protocol))
	}
	return strings.Join(ports, ", ")
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
