package plugin

import (
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Runners is an interface that manager can call to get runners for a
// particular action.
type Runners interface {
	// Print returns a runner for printing. The caller should close
	// the channel when they are done with it.
	Print(ManagerStore) (DefaultRunner, chan PrintResponse)
	// Tab returns a runner for tabs. The caller should close
	// the channel when they are done with it.
	Tab(ManagerStore) (DefaultRunner, chan component.Tab)
}

type defaultRunners struct{}

var _ Runners = (*defaultRunners)(nil)

func newDefaultRunners() *defaultRunners {
	return &defaultRunners{}
}

func (dr *defaultRunners) Print(store ManagerStore) (DefaultRunner, chan PrintResponse) {
	ch := make(chan PrintResponse)
	return PrintRunner(store, ch), ch
}

func (dr *defaultRunners) Tab(store ManagerStore) (DefaultRunner, chan component.Tab) {
	ch := make(chan component.Tab)
	return TabRunner(store, ch), ch
}

// DefaultRunner runs a function against all plugins
type DefaultRunner struct {
	RunFunc func(name string, gvk schema.GroupVersionKind, object runtime.Object) error
}

// Run runs the runner for an object with the provided clients.
func (pr *DefaultRunner) Run(object runtime.Object, clientNames []string) error {
	if err := pr.validate(object); err != nil {
		return errors.Wrap(err, "plugin runner validate")
	}

	var g errgroup.Group

	gvk := object.GetObjectKind().GroupVersionKind()

	for _, name := range clientNames {
		fn := func(name string) func() error {
			return func() error {
				if err := pr.RunFunc(name, gvk, object); err != nil {
					return errors.Wrap(err, "running")
				}

				return nil
			}
		}

		g.Go(fn(name))
	}

	if err := g.Wait(); err != nil {
		return errors.Wrap(err, "handle object")
	}

	return nil
}

func (pr *DefaultRunner) validate(object runtime.Object) error {
	if object == nil {
		return errors.New("object is nil")
	}

	if pr.RunFunc == nil {
		return errors.New("requires a runFunc")
	}

	return nil
}

// PrintRunner is a runner for printing.
func PrintRunner(store ManagerStore, ch chan<- PrintResponse) DefaultRunner {
	return DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			metadata, err := store.GetMetadata(name)
			if err != nil {
				return err
			}

			if !metadata.Capabilities.HasPrinterSupport(gvk) {
				return nil
			}

			resp, err := printObject(store, name, object)
			if err != nil {
				return err
			}

			ch <- resp
			return nil
		},
	}
}

func printObject(store ManagerStore, pluginName string, object runtime.Object) (PrintResponse, error) {
	if store == nil {
		return PrintResponse{}, errors.New("store is nil")
	}

	service, err := store.GetService(pluginName)
	if err != nil {
		return PrintResponse{}, err
	}

	resp, err := service.Print(object)
	if err != nil {
		return PrintResponse{}, errors.Wrapf(err, "print object with plugin %q", pluginName)
	}

	return resp, nil
}

// TabRunner is a runner for tabs.
func TabRunner(store ManagerStore, ch chan<- component.Tab) DefaultRunner {
	runner := DefaultRunner{
		RunFunc: func(name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			if store == nil {
				return errors.New("store is nil")
			}

			metadata, err := store.GetMetadata(name)
			if err != nil {
				return err
			}

			if !metadata.Capabilities.HasTabSupport(gvk) {
				return nil
			}

			service, err := store.GetService(name)
			if err != nil {
				return err
			}

			tab, err := service.PrintTab(object)
			if err != nil {
				return errors.Wrapf(err, "printing tab for plugin %q", name)
			}

			ch <- *tab

			return nil
		},
	}

	return runner
}
