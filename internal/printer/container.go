/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"

	"path"
	"sort"
	"strings"

	"github.com/vmware-tanzu/octant/internal/octant"

	"github.com/vmware-tanzu/octant/internal/portforward"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

type containerActionFunc func(runtime.Object, *corev1.Container) (component.Action, error)

// ContainerConfigurationOption
type ContainerConfigurationOption func(*ContainerConfiguration)

func IsInit(b bool) ContainerConfigurationOption {
	return func(cc *ContainerConfiguration) {
		cc.isInit = b
	}
}

func WithPrintOptions(options Options) ContainerConfigurationOption {
	return func(cc *ContainerConfiguration) {
		cc.options = options
	}
}

func WithActions(actions ...containerActionFunc) ContainerConfigurationOption {
	return func(cc *ContainerConfiguration) {
		for _, action := range actions {
			cc.actionGenerators = append(cc.actionGenerators, action)
		}
	}
}

// ContainerConfiguration generates container configuration.
type ContainerConfiguration struct {
	parent             runtime.Object
	container          *corev1.Container
	portForwardService portforward.PortForwarder
	isInit             bool
	options            Options
	context            context.Context
	actionGenerators   []containerActionFunc
}

// NewContainerConfiguration creates an instance of ContainerConfiguration.
func NewContainerConfiguration(ctx context.Context, parent runtime.Object, c *corev1.Container, pfs portforward.PortForwarder, opts ...ContainerConfigurationOption) *ContainerConfiguration {
	cc := &ContainerConfiguration{
		parent:             parent,
		container:          c,
		portForwardService: pfs,
		context:            ctx,
		actionGenerators:   []containerActionFunc{},
	}
	for _, o := range opts {
		o(cc)
	}
	return cc
}

// Create creates a deployment configuration summary.
func (cc *ContainerConfiguration) Create() (*component.Summary, error) {
	if cc == nil || cc.container == nil {
		return nil, errors.New("container is nil")
	}
	c := cc.container

	var containerStatus *corev1.ContainerStatus
	if pod, ok := cc.parent.(*corev1.Pod); ok {
		var err error
		containerStatus, err = findContainerStatus(pod, cc.container.Name, cc.isInit)
		if err != nil {
			switch err.(type) {
			case *containerNotFoundError:
			// no op
			default:
				return nil, errors.Wrapf(err, "get container status for %q", cc.container.Name)
			}
		}
	}

	sections := component.SummarySections{}

	sections.AddText("Image", c.Image)
	if containerStatus != nil {
		sections.AddText("Image ID", containerStatus.ImageID)
	}

	hostPorts := describeContainerHostPorts(c.Ports)
	if hostPorts != "" {
		sections.AddText("Host Ports", hostPorts)
	}
	containerPorts, err := describeContainerPorts(cc.context, cc.parent, c.Ports, cc.portForwardService)
	if err != nil {
		return nil, errors.Wrap(err, "describe container ports")
	}
	if len(containerPorts) > 0 {
		sections.Add("Container Ports", component.NewPorts(containerPorts))
	}

	if containerStatus != nil && !cc.isInit {
		lastState, found := printContainerState(containerStatus.LastTerminationState)
		if found {
			sections.AddText("Last State", lastState)
		}
		currentState, found := printContainerState(containerStatus.State)
		if found {
			sections.AddText("Current State", currentState)
		}

		sections.AddText("Ready", fmt.Sprintf("%t", containerStatus.Ready))
		sections.AddText("Restart Count", fmt.Sprintf("%d", containerStatus.RestartCount))
	}

	envTbl, err := describeContainerEnv(cc.context, cc.parent, c, cc.options)
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

	for _, actionFunc := range cc.actionGenerators {
		action, err := actionFunc(cc.parent, c)
		if err != nil {
			logger := log.From(cc.context)
			logger.Errorf("failed loading actions: %w", err)
			continue
		}
		summary.AddAction(action)
	}

	return summary, nil
}

func printContainerState(state corev1.ContainerState) (string, bool) {
	switch {
	case state.Running != nil:
		return fmt.Sprintf("started at %s", state.Running.StartedAt), true
	case state.Waiting != nil:
		return fmt.Sprintf("waiting: %s", state.Waiting.Message), true
	case state.Terminated != nil:
		return fmt.Sprintf("terminated with %d at %s: %s",
			state.Terminated.ExitCode,
			state.Terminated.FinishedAt,
			state.Terminated.Reason), true
	}

	return "indeterminate", false
}

type containerStatus interface {
	isContainerFound() bool
}

type containerNotFoundError struct {
	name string
}

var _ containerStatus = (*containerNotFoundError)(nil)

func (e *containerNotFoundError) Error() string {
	return fmt.Sprintf("container %q not found", e.name)
}

func (e *containerNotFoundError) isContainerFound() bool {
	return false
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
	ctx context.Context,
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

	states, err := portForwardService.FindPod(namespace, gvk, name)
	if err != nil {
		if _, ok := err.(notFound); !ok {
			return nil, errors.Wrap(err, "query port forward service for pod")
		}
	}

	pfLookup := make(map[int32]string, len(states))
	for _, cPort := range cPorts {
		for _, state := range states {
			for _, forwarded := range state.Ports {
				if int(forwarded.Remote) == int(cPort.ContainerPort) {
					pfLookup[cPort.ContainerPort] = state.ID
				}
			}
		}

		if cPort.ContainerPort == 0 {
			continue
		}
		pfs := component.PortForwardState{}
		var port *component.Port

		if isPod && cPort.Protocol == corev1.ProtocolTCP {
			pfs.IsForwardable = true
		}

		if id, ok := pfLookup[cPort.ContainerPort]; ok {
			if state, found := portForwardService.Get(id); found {
				for _, forwarded := range state.Ports {
					pfs.Port = int(forwarded.Local)
				}
			} else {
				return nil, errors.Wrap(err, "id not found for port forward")
			}
			pfs.ID = id
			pfs.IsForwarded = true
		}

		apiVersion, kind := gvk.ToAPIVersionAndKind()
		port = component.NewPort(namespace, apiVersion, kind, name, int(cPort.ContainerPort), string(cPort.Protocol), pfs)
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
func describeContainerEnv(ctx context.Context, parent runtime.Object, c *corev1.Container, options Options) (*component.Table, error) {
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
	tbl := component.NewTable("Environment", "There are no defined environment variables!", cols)

	envRows, err := describeEnvRows(ctx, ns, c.Env, options)
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
func describeEnvRows(ctx context.Context, namespace string, vars []corev1.EnvVar, options Options) ([]component.TableRow, error) {
	rows := make([]component.TableRow, 0)

	sort.Slice(vars, func(i, j int) bool { return vars[i].Name < vars[j].Name })

	for _, e := range vars {
		row := component.TableRow{}
		rows = append(rows, row)

		// each row requires data for each defined column
		row["Source"] = component.NewText("")

		row["Name"] = component.NewText(e.Name)
		row["Value"] = component.NewText(e.Value)
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

			objectStore := options.DashConfig.ObjectStore()

			key := store.Key{
				Name:       ref.Name,
				Namespace:  namespace,
				APIVersion: "v1",
				Kind:       "ConfigMap",
			}

			u, err := objectStore.Get(ctx, key)
			if err != nil {
				return nil, err
			}

			if u != nil {
				configMap := &corev1.ConfigMap{}
				if err := kubernetes.FromUnstructured(u, configMap); err != nil {
					return nil, err
				}

				row["Value"] = component.NewText(configMap.Data[ref.Key])
				row["Source"] = source
			} else {
				row["Value"] = component.NewText("<none>")
				row["Source"] = component.NewText(fmt.Sprintf("%s:%s", ref.Name, ref.Key))
			}
		}
	}

	return rows, nil
}

// describeEnvFromRows renders container environmentFrom references as table rows.
// Expected columns: Name, Value, Source
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
	tbl := component.NewTable("Volume Mounts", "There are no volume mounts!", cols)

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

func editContainerAction(owner runtime.Object, container *corev1.Container) (component.Action, error) {
	if container == nil {
		return component.Action{}, errors.New("container is nil")
	}

	containersPath, err := containersPathForObject(owner)
	if err != nil {
		return component.Action{}, err
	}

	containersPathData, err := json.Marshal(containersPath)
	if err != nil {
		return component.Action{}, err
	}

	form, err := component.CreateFormForObject(octant.ActionOverviewContainerEditor, owner,
		component.NewFormFieldText("Image", "containerImage", container.Image),
		component.NewFormFieldHidden("containersPath", string(containersPathData)),
		component.NewFormFieldHidden("containerName", container.Name),
	)
	if err != nil {
		return component.Action{}, err
	}

	action := component.Action{
		Name:  "Edit",
		Title: fmt.Sprintf("Container %s Editor", container.Name),
		Form:  form,
	}

	return action, nil
}

func containersPathForObject(object runtime.Object) ([]string, error) {
	if object == nil {
		return nil, errors.New("object is nil")
	}

	g := object.GetObjectKind().GroupVersionKind()

	switch {
	case g.Group == "batch" && g.Kind == "CronJob":
		return []string{"spec", "jobTemplate", "spec", "template", "containers"}, nil
	case g.Group == "batch" && g.Kind == "Job":
		return []string{"spec", "template", "containers"}, nil
	case g.Group == "apps" && g.Kind == "DaemonSet":
		return []string{"spec", "template", "spec", "containers"}, nil
	case g.Group == "apps" && g.Kind == "Deployment":
		return []string{"spec", "template", "spec", "containers"}, nil
	case g.Group == "apps" && g.Kind == "ReplicaSet":
		return []string{"spec", "template", "spec", "containers"}, nil
	case g.Group == "" && g.Kind == "ReplicationController":
		return []string{"spec", "template", "spec", "containers"}, nil
	case g.Group == "apps" && g.Kind == "StatefulSet":
		return []string{"spec", "template", "spec", "containers"}, nil
	default:
		return nil, errors.Errorf("unable to find containers location for %+v, %s, %s", g, g.Group, g.Kind)
	}
}
