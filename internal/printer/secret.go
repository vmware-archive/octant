/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	secretTableCols = component.NewTableCols("Name", "Labels", "Type", "Data", "Age")
	secretDataCols  = component.NewTableCols("Key")
)

// SecretListHandler is a printFunc that lists secrets.
func SecretListHandler(ctx context.Context, list *corev1.SecretList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("list of secrets is nil")
	}

	ot := NewObjectTable("Secrets", "We couldn't find any secrets!", secretTableCols, options.DashConfig.ObjectStore())

	for _, secret := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&secret, secret.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

		row["Labels"] = component.NewLabels(secret.ObjectMeta.Labels)
		row["Type"] = component.NewText(string(secret.Type))
		row["Data"] = component.NewText(fmt.Sprintf("%d", len(secret.Data)))
		row["Age"] = component.NewTimestamp(secret.ObjectMeta.CreationTimestamp.Time)

		if err := ot.AddRowForObject(ctx, &secret, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// SecretHandler is a printFunc for printing a secret summary.
func SecretHandler(ctx context.Context, secret *corev1.Secret, options Options) (component.Component, error) {
	o := NewObject(secret)

	sh, err := newSecretHandler(secret, o)
	if err != nil {
		return nil, err
	}

	if err := sh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print secret configuration")
	}

	if err := sh.Data(options); err != nil {
		return nil, errors.Wrap(err, "print secret data")
	}

	return o.ToComponent(ctx, options)
}

// SecretConfiguration generates a secret configuration
type SecretConfiguration struct {
	secret *corev1.Secret
}

// NewSecretConfiguration creates an instance of SecretConfiguration
func NewSecretConfiguration(secret *corev1.Secret) *SecretConfiguration {
	return &SecretConfiguration{
		secret: secret,
	}
}

// Create creates a secret configuration summary
func (s *SecretConfiguration) Create(options Options) (*component.Summary, error) {
	if s.secret == nil {
		return nil, errors.New("secret is nil")
	}
	secret := s.secret

	var sections []component.SummarySection

	sections = append(sections, component.SummarySection{
		Header:  "Type",
		Content: component.NewText(string(secret.Type)),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func describeSecretData(secret corev1.Secret) (*component.Table, error) {
	table := component.NewTable("Data", "This secret has no data!", secretDataCols)

	for key := range secret.Data {
		row := component.TableRow{}
		row["Key"] = component.NewText(key)

		table.Add(row)
	}

	table.Sort(false, "Key")

	return table, nil
}

type secretObject interface {
	Config(options Options) error
	Data(options Options) error
}

type secretHandler struct {
	secret     *corev1.Secret
	configFunc func(*corev1.Secret, Options) (*component.Summary, error)
	dataFunc   func(*corev1.Secret, Options) (*component.Table, error)
	object     *Object
}

func newSecretHandler(secret *corev1.Secret, object *Object) (*secretHandler, error) {
	if secret == nil {
		return nil, errors.New("can't print a nil secret")
	}

	if object == nil {
		return nil, errors.New("can't print secret using a nil object printer")
	}

	sh := &secretHandler{
		secret:     secret,
		configFunc: defaultSecretConfig,
		dataFunc:   defaultSecretData,
		object:     object,
	}

	return sh, nil
}

func (s *secretHandler) Config(options Options) error {
	out, err := s.configFunc(s.secret, options)
	if err != nil {
		return err
	}
	s.object.RegisterConfig(out)
	return nil
}

func defaultSecretConfig(secret *corev1.Secret, options Options) (*component.Summary, error) {
	return NewSecretConfiguration(secret).Create(options)
}

func (s *secretHandler) Data(options Options) error {
	if s.secret == nil {
		return errors.New("can't display data for nil secret")
	}

	s.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return s.dataFunc(s.secret, options)
		},
	})
	return nil
}

func defaultSecretData(secret *corev1.Secret, options Options) (*component.Table, error) {
	return describeSecretData(*secret)
}
