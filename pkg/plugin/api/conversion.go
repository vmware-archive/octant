/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/vmware-tanzu/octant/internal/util/json"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/plugin/api/proto"
	"github.com/vmware-tanzu/octant/pkg/store"
)

func convertFromKey(in store.Key) (*proto.KeyRequest, error) {
	keyRequest := proto.KeyRequest{
		Namespace:  in.Namespace,
		ApiVersion: in.APIVersion,
		Kind:       in.Kind,
		Name:       in.Name,
	}

	if in.Selector != nil {
		keyRequest.LabelSelector = &wrapperspb.BytesValue{Value: []byte(in.Selector.String())}
	}

	return &keyRequest, nil
}

func convertToKey(in *proto.KeyRequest) (store.Key, error) {
	if in == nil {
		return store.Key{}, errors.New("key request is nil")
	}

	key := store.Key{
		Namespace:  in.Namespace,
		APIVersion: in.ApiVersion,
		Kind:       in.Kind,
		Name:       in.Name,
	}

	labelSelector := in.GetLabelSelector()
	if labelSelector != nil {
		selector, err := metav1.ParseToLabelSelector(string(labelSelector.Value))
		if err != nil {
			return store.Key{}, errors.New("cannot parse selector string")
		}

		matchLabels := labels.Set{}
		for label, value := range selector.MatchLabels {
			matchLabels[label] = value
		}
		key.Selector = &matchLabels
	}

	return key, nil
}

func convertFromAlert(alert action.Alert) (*proto.AlertRequest, error) {
	if alert.Expiration == nil {
		return &proto.AlertRequest{}, errors.New("expiration is nil")
	}

	alertRequest := proto.AlertRequest{
		Type:       string(alert.Type),
		Message:    alert.Message,
		Expiration: timestamppb.New(*alert.Expiration),
	}

	return &alertRequest, nil
}

func convertToAlert(in *proto.AlertRequest) (action.Alert, error) {
	if in == nil {
		return action.Alert{}, errors.New("alert request is nil")
	}

	expiration := in.Expiration.AsTime()

	alert := action.Alert{
		Type:       action.AlertType(in.Type),
		Message:    in.Message,
		Expiration: &expiration,
	}

	return alert, nil
}

func convertToLinkComponent(in, name string) (*component.Link, error) {
	if in == "" {
		return &component.Link{}, nil
	}
	link := component.NewLink("", name, in)
	return link, nil
}

func convertFromPayload(payload action.Payload) ([]byte, error) {
	return json.Marshal(payload)
}

func convertToEvent(in *proto.EventRequest) (string, string, action.Payload, error) {
	if in == nil {
		return "", "", action.Payload{}, fmt.Errorf("event request is nil")
	}
	payload := action.Payload{}
	err := json.Unmarshal(in.Payload, &payload)
	if err != nil {
		return "", "", action.Payload{}, fmt.Errorf("unmarshal payload: %w", err)
	}
	return in.ClientID, in.EventName, payload, nil
}

func convertFromObjects(in *unstructured.UnstructuredList) ([][]byte, error) {
	var out [][]byte

	for _, object := range in.Items {
		data, err := convertFromObject(&object)
		if err != nil {
			return nil, err
		}

		out = append(out, data)
	}

	return out, nil
}

func convertFromObject(in *unstructured.Unstructured) ([]byte, error) {
	return json.Marshal(in)
}

func convertToObjects(in [][]byte) (*unstructured.UnstructuredList, error) {
	list := &unstructured.UnstructuredList{}

	for _, data := range in {
		object, err := convertToObject(data)
		if err != nil {
			return nil, err
		}
		if object == nil {
			continue
		}
		list.Items = append(list.Items, *object)
	}

	return list, nil
}

func convertToObject(in []byte) (*unstructured.Unstructured, error) {
	if in == nil {
		return nil, nil
	}

	object := unstructured.Unstructured{}
	err := json.Unmarshal(in, &object)
	if err != nil {
		return nil, err
	}

	if object.Object == nil {
		return nil, nil
	}

	return &object, nil
}

func convertToPortForwardRequest(in *proto.PortForwardRequest) (*PortForwardRequest, error) {
	if in == nil {
		return nil, errors.New("can't convert nil object")
	}

	port := in.PortNumber
	if port > 0xFFFF {
		return nil, errors.Errorf("port number must be a uint32; it was: %d", port)
	}

	return &PortForwardRequest{
		Namespace: in.Namespace,
		PodName:   in.PodName,
		Port:      uint16(port),
	}, nil
}
