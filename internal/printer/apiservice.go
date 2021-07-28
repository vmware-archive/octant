/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

// APIServiceListHandler is a printFunc that prints api services
func APIServiceListHandler(ctx context.Context, list *apiregistrationv1.APIServiceList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("api service list is nil")
	}

	cols := component.NewTableCols("Name", "Service", "Age")
	ot := NewObjectTable("API Services", "We couldn't find any api services!", cols, options.DashConfig.ObjectStore())
	ot.EnablePluginStatus(options.DashConfig.PluginManager())

	for _, apiService := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&apiService, apiService.Name)
		if err != nil {
			return nil, err
		}

		var service component.Component
		if apiService.Spec.Service == nil {
			service = component.NewText("Local")
		} else {
			serviceLink, err := options.Link.ForGVK(
				apiService.Spec.Service.Namespace,
				gvk.Service.GroupVersion().String(),
				gvk.Service.Kind, apiService.Spec.Service.Name,
				fmt.Sprintf("%s/%s", apiService.Spec.Service.Namespace, apiService.Spec.Service.Name),
			)
			if err != nil {
				return nil, err
			}
			service = serviceLink
		}

		row["Name"] = nameLink
		row["Service"] = service
		ts := apiService.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		if err := ot.AddRowForObject(ctx, &apiService, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// APIServiceHandler is a printFunc that prints a api service
func APIServiceHandler(ctx context.Context, apiService *apiregistrationv1.APIService, options Options) (component.Component, error) {
	o := NewObject(apiService)

	ch, err := newAPIServiceHandler(apiService, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Config(options); err != nil {
		return nil, errors.Wrap(err, "print apiservice configuration")
	}

	if err := ch.Status(options); err != nil {
		return nil, errors.Wrap(err, "print apiservice status")
	}

	return o.ToComponent(ctx, options)
}

// APIServiceConfiguration generates a apiservice configuration
type APIServiceConfiguration struct {
	apiService *apiregistrationv1.APIService
}

// NewAPIServiceConfiguration creates an instance of APIServiceConfiguration
func NewAPIServiceConfiguration(apiService *apiregistrationv1.APIService) *APIServiceConfiguration {
	return &APIServiceConfiguration{
		apiService: apiService,
	}
}

// Create creates a apiservice configuration summary
func (c *APIServiceConfiguration) Create(options Options) (*component.Summary, error) {
	if c == nil || c.apiService == nil {
		return nil, errors.New("apiservice is nil")
	}

	var sections component.SummarySections

	service, err := apiServiceService(c.apiService, options)
	if err != nil {
		return nil, err
	}
	sections.Add("Service", service)
	tls := apiServiceTLS(c.apiService)
	sections.AddText("TLS", tls)
	sections.AddText("Group", c.apiService.Spec.Group)
	groupPriority := fmt.Sprintf("%d", c.apiService.Spec.GroupPriorityMinimum)
	sections.AddText("Group Priority Minimum", groupPriority)
	sections.AddText("Version", c.apiService.Spec.Version)
	versionPriority := fmt.Sprintf("%d", c.apiService.Spec.VersionPriority)
	sections.AddText("Version Priority", versionPriority)

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

// APIServiceStatus generates a apiservice status
type APIServiceStatus struct {
	apiService *apiregistrationv1.APIService
}

// NewAPIServiceStatus creates an instance of APIServiceStatus
func NewAPIServiceStatus(apiService *apiregistrationv1.APIService) *APIServiceStatus {
	return &APIServiceStatus{
		apiService: apiService,
	}
}

// Create creates a apiservice status summary
func (c *APIServiceStatus) Create(options Options) (*component.Summary, error) {
	if c.apiService == nil {
		return nil, errors.New("apiService is nil")
	}

	summary := component.NewSummary("Status")

	sections := component.SummarySections{}

	var availableCond *apiregistrationv1.APIServiceCondition
	for _, cond := range c.apiService.Status.Conditions {
		if cond.Type == apiregistrationv1.Available {
			availableCond = &cond
			break
		}
	}
	if availableCond == nil {
		availableCond = &apiregistrationv1.APIServiceCondition{}
	}

	if availableCond.Status == "" {
		sections.AddText("Available", "Unknown")
	} else {
		sections.AddText("Available", string(availableCond.Status))
	}
	if availableCond.Reason != "" {
		sections.AddText("Reason", availableCond.Reason)
	}
	if availableCond.Message != "" {
		sections.AddText("Message", availableCond.Message)
	}

	summary.Add(sections...)

	return summary, nil
}

type apiServiceObject interface {
	Config(options Options) error
	Status(options Options) error
}

type apiServiceHandler struct {
	apiService  *apiregistrationv1.APIService
	configFunc  func(*apiregistrationv1.APIService, Options) (*component.Summary, error)
	summaryFunc func(*apiregistrationv1.APIService, Options) (*component.Summary, error)
	object      *Object
}

var _ apiServiceObject = (*apiServiceHandler)(nil)

func newAPIServiceHandler(apiService *apiregistrationv1.APIService, object *Object) (*apiServiceHandler, error) {
	if apiService == nil {
		return nil, errors.New("can't print a nil apiservice")
	}

	if object == nil {
		return nil, errors.New("can't print a apiservice using an nil object printer")
	}

	ch := &apiServiceHandler{
		apiService:  apiService,
		configFunc:  defaultAPIServiceConfig,
		summaryFunc: defaultAPIServiceSummary,
		object:      object,
	}
	return ch, nil
}

func (c *apiServiceHandler) Config(options Options) error {
	out, err := c.configFunc(c.apiService, options)
	if err != nil {
		return err
	}
	c.object.RegisterConfig(out)
	return nil
}

func (c *apiServiceHandler) Status(options Options) error {
	out, err := c.summaryFunc(c.apiService, options)
	if err != nil {
		return err
	}
	c.object.RegisterSummary(out)
	return nil
}

func defaultAPIServiceConfig(apiService *apiregistrationv1.APIService, options Options) (*component.Summary, error) {
	return NewAPIServiceConfiguration(apiService).Create(options)
}

func defaultAPIServiceSummary(apiService *apiregistrationv1.APIService, options Options) (*component.Summary, error) {
	return NewAPIServiceStatus(apiService).Create(options)
}

func apiServiceService(apiService *apiregistrationv1.APIService, options Options) (component.Component, error) {
	if apiService.Spec.Service == nil {
		return component.NewText("Local"), nil
	}
	return options.Link.ForGVK(
		apiService.Spec.Service.Namespace,
		gvk.Service.GroupVersion().String(),
		gvk.Service.Kind, apiService.Spec.Service.Name,
		fmt.Sprintf("%s/%s", apiService.Spec.Service.Namespace, apiService.Spec.Service.Name),
	)
}

func apiServiceTLS(apiService *apiregistrationv1.APIService) string {
	if apiService.Spec.InsecureSkipTLSVerify {
		return "Skip verification"
	}
	if len(apiService.Spec.CABundle) != 0 {
		return "Custom CA"
	}
	return "System trust store"
}
