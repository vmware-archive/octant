/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware-tanzu/octant/pkg/view/component"
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
	// ObjectStatus returns a runner for object status. The caller should
	// close the channel when they are done with it.
	ObjectStatus(ManagerStore) (DefaultRunner, chan ObjectStatusResponse)
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

func (dr *defaultRunners) ObjectStatus(store ManagerStore) (DefaultRunner, chan ObjectStatusResponse) {
	ch := make(chan ObjectStatusResponse)
	return ObjectStatusRunner(store, ch), ch
}

// DefaultRunner runs a function against all plugins
type DefaultRunner struct {
	RunFunc func(ctx context.Context, name string, gvk schema.GroupVersionKind, object runtime.Object) error
}

// Run runs the runner for an object with the provided clients.
func (pr *DefaultRunner) Run(ctx context.Context, object runtime.Object, clientNames []string) error {
	if err := pr.validate(object); err != nil {
		return fmt.Errorf("plugin runner validate: %w", err)
	}

	var g errgroup.Group

	gvk := object.GetObjectKind().GroupVersionKind()

	for _, name := range clientNames {
		fn := func(name string) func() error {
			return func() error {
				if err := pr.RunFunc(ctx, name, gvk, object); err != nil {
					return fmt.Errorf("running on %s: %w", name, err)
				}

				return nil
			}
		}
		g.Go(fn(name))
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("handle object: %w", err)
	}

	return nil
}

func (pr *DefaultRunner) validate(object runtime.Object) error {
	if object == nil {
		return fmt.Errorf("object is nil")
	}

	if pr.RunFunc == nil {
		return fmt.Errorf("requires a runFunc")
	}

	return nil
}

// PrintRunner is a runner for printing.
func PrintRunner(store ManagerStore, ch chan<- PrintResponse) DefaultRunner {
	return DefaultRunner{
		RunFunc: func(ctx context.Context, name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			if IsJavaScriptPlugin(name) {
				jsPlugin, ok := store.GetJS(name)
				if !ok {
					return fmt.Errorf("plugin %s not found", name)
				}
				if !jsPlugin.Metadata().Capabilities.HasPrinterSupport(gvk) {
					return nil
				}

				resp, err := jsPlugin.Print(ctx, object)

				if err != nil {
					return err
				}
				ch <- resp
				return nil
			}

			metadata, err := store.GetMetadata(name)
			if err != nil {
				return err
			}

			if !metadata.Capabilities.HasPrinterSupport(gvk) {
				return nil
			}

			resp, err := printObject(ctx, store, name, object)
			if err != nil {
				return err
			}

			ch <- resp
			return nil
		},
	}
}

func printObject(ctx context.Context, store ManagerStore, pluginName string, object runtime.Object) (PrintResponse, error) {
	if store == nil {
		return PrintResponse{}, fmt.Errorf("store is nil")
	}

	service, err := store.GetService(pluginName)
	if err != nil {
		return PrintResponse{}, err
	}

	resp, err := service.Print(ctx, object)
	if err != nil {
		return PrintResponse{}, fmt.Errorf("print object with plugin %q: %w", pluginName, err)
	}

	return resp, nil
}

// TabRunner is a runner for tabs.
func TabRunner(store ManagerStore, ch chan<- component.Tab) DefaultRunner {
	runner := DefaultRunner{
		RunFunc: func(ctx context.Context, name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			if IsJavaScriptPlugin(name) {
				jsPlugin, ok := store.GetJS(name)
				if !ok {
					return fmt.Errorf("plugin %s not found", name)
				}

				if !jsPlugin.Metadata().Capabilities.HasTabSupport(gvk) {
					return nil
				}

				resp, err := jsPlugin.PrintTab(ctx, object)
				if err != nil {
					return fmt.Errorf("printing tabResponse for plugin: %q: %w", name, err)
				}

				ch <- *resp.Tab
				return nil
			}

			if store == nil {
				return fmt.Errorf("store is nil")
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

			tabResponse, err := service.PrintTab(ctx, object)
			if err != nil {
				return fmt.Errorf("printing tabResponse for plugin %q: %w", name, err)
			}

			ch <- *tabResponse.Tab

			return nil
		},
	}

	return runner
}

// ObjectStatusRunner is a runner for object status.
func ObjectStatusRunner(store ManagerStore, ch chan<- ObjectStatusResponse) DefaultRunner {
	return DefaultRunner{
		RunFunc: func(ctx context.Context, name string, gvk schema.GroupVersionKind, object runtime.Object) error {
			if IsJavaScriptPlugin(name) {
				jsPlugin, ok := store.GetJS(name)
				if !ok {
					return fmt.Errorf("plugin %s not found", name)
				}

				if !jsPlugin.Metadata().Capabilities.HasObjectStatusSupport(gvk) {
					return nil
				}

				resp, err := jsPlugin.ObjectStatus(ctx, object)
				if err != nil {
					return fmt.Errorf("printing objectStatus for plugin: %q: %w", name, err)
				}

				ch <- resp
				return nil
			}

			metadata, err := store.GetMetadata(name)
			if err != nil {
				return err
			}

			if !metadata.Capabilities.HasObjectStatusSupport(gvk) {
				return nil
			}

			service, err := store.GetService(name)
			if err != nil {
				return err
			}

			resp, err := service.ObjectStatus(ctx, object)
			if err != nil {
				return fmt.Errorf("print object status with plugin %q: %w", name, err)
			}

			ch <- resp
			return nil
		},
	}
}
