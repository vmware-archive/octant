/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
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
		return errors.Wrap(err, "plugin runner validate")
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
		RunFunc: func(ctx context.Context, name string, gvk schema.GroupVersionKind, object runtime.Object) error {
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
		return PrintResponse{}, errors.New("store is nil")
	}

	service, err := store.GetService(pluginName)
	if err != nil {
		return PrintResponse{}, err
	}

	resp, err := service.Print(ctx, object)
	if err != nil {
		return PrintResponse{}, errors.Wrapf(err, "print object with plugin %q", pluginName)
	}

	return resp, nil
}

// TabRunner is a runner for tabs.
func TabRunner(store ManagerStore, ch chan<- component.Tab) DefaultRunner {
	runner := DefaultRunner{
		RunFunc: func(ctx context.Context, name string, gvk schema.GroupVersionKind, object runtime.Object) error {
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

			tabResponse, err := service.PrintTab(ctx, object)
			if err != nil {
				return errors.Wrapf(err, "printing tabResponse for plugin %q", name)
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
			metadata, err := store.GetMetadata(name)
			if err != nil {
				return err
			}

			if !metadata.Capabilities.HasPrinterSupport(gvk) {
				return nil
			}

			service, err := store.GetService(name)
			if err != nil {
				return err
			}

			resp, err := service.ObjectStatus(ctx, object)
			if err != nil {
				return errors.Wrapf(err, "print object status with plugin %q", name)
			}

			ch <- resp
			return nil
		},
	}
}
