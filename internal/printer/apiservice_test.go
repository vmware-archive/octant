/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	apiregistrationv1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_APIServiceListHandler(t *testing.T) {
	object := testutil.CreateAPIService("v1", "example.com")
	object.CreationTimestamp = *testutil.CreateTimestamp()
	object.Spec.Service = &apiregistrationv1.ServiceReference{
		Namespace: "default",
		Name:      "service",
	}

	list := &apiregistrationv1.APIServiceList{
		Items: []apiregistrationv1.APIService{*object},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(object, object.Name, "/path")
	tpo.PathForGVK("default", "v1", "Service", "service", "default/service", "/service")

	now := testutil.Time()

	ctx := context.Background()
	got, err := APIServiceListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Service", "Age")
	expected := component.NewTable("API Services", "We couldn't find any api services!", cols)

	expected.Add(component.TableRow{
		"Name": component.NewLink("", object.Name, "/path",
			genObjectStatus(component.TextStatusOK, []string{
				"apiregistration.k8s.io/v1 APIService is OK",
			})),
		"Service": component.NewLink("", "default/service", "/service"),
		"Age":     component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	testutil.AssertJSONEqual(t, expected, got)
}

func Test_APIServiceListHandler_local(t *testing.T) {
	object := testutil.CreateAPIService("v1", "apps")
	object.CreationTimestamp = *testutil.CreateTimestamp()

	list := &apiregistrationv1.APIServiceList{
		Items: []apiregistrationv1.APIService{*object},
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(object, object.Name, "/path")

	now := testutil.Time()

	ctx := context.Background()
	got, err := APIServiceListHandler(ctx, list, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Service", "Age")
	expected := component.NewTable("API Services", "We couldn't find any api services!", cols)

	expected.Add(component.TableRow{
		"Name": component.NewLink("", object.Name, "/path",
			genObjectStatus(component.TextStatusOK, []string{
				"apiregistration.k8s.io/v1 APIService is OK",
			})),
		"Service": component.NewText("Local"),
		"Age":     component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, object),
		}),
	})

	testutil.AssertJSONEqual(t, expected, got)
}

func Test_APIServerConfiguration(t *testing.T) {
	apiService := testutil.CreateAPIService("v1", "example.com")

	cases := []struct {
		name       string
		apiService *apiregistrationv1.APIService
		isErr      bool
		expected   *component.Summary
	}{
		{
			name:       "general",
			apiService: apiService,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Service",
					Content: component.NewText("Local"),
				},
				{
					Header:  "TLS",
					Content: component.NewText("System trust store"),
				},
				{
					Header:  "Group",
					Content: component.NewText("example.com"),
				},
				{
					Header:  "Group Priority Minimum",
					Content: component.NewText("100"),
				},
				{
					Header:  "Version",
					Content: component.NewText("v1"),
				},
				{
					Header:  "Version Priority",
					Content: component.NewText("100"),
				},
			}...),
		},
		{
			name:       "apiservice is nil",
			apiService: nil,
			isErr:      true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			cc := NewAPIServiceConfiguration(test.apiService)

			summary, err := cc.Create(printOptions)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, test.expected, summary)
		})
	}
}

func Test_APIServerStatus(t *testing.T) {
	apiService := testutil.CreateAPIService("v1", "example.com")
	apiServiceUnknown := apiService.DeepCopy()
	apiServiceUnknown.Status.Conditions = nil

	cases := []struct {
		name       string
		apiService *apiregistrationv1.APIService
		isErr      bool
		expected   *component.Summary
	}{
		{
			name:       "general",
			apiService: apiService,
			expected: component.NewSummary("Status", []component.SummarySection{
				{
					Header:  "Available",
					Content: component.NewText("True"),
				},
				{
					Header:  "Reason",
					Content: component.NewText("Local"),
				},
				{
					Header:  "Message",
					Content: component.NewText("Local APIServices are always available"),
				},
			}...),
		},
		{
			name:       "unknown",
			apiService: apiServiceUnknown,
			expected: component.NewSummary("Status", []component.SummarySection{
				{
					Header:  "Available",
					Content: component.NewText("Unknown"),
				},
			}...),
		},
		{
			name:       "apiservice is nil",
			apiService: nil,
			isErr:      true,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			cc := NewAPIServiceStatus(test.apiService)

			summary, err := cc.Create(printOptions)
			if test.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, test.expected, summary)
		})
	}
}
