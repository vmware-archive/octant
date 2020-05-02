/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vmware-tanzu/octant/pkg/log"
)

//go:generate mockgen -destination=./fake/mock_alert.go -package=fake github.com/vmware-tanzu/octant/pkg/action Alerter

const (
	// DefaultAlertExpiration is the default expiration for alerts.
	DefaultAlertExpiration = 10 * time.Second
)

// AlertType is the type of alert.
type AlertType string

const (
	// AlertTypeError is for error alerts.
	AlertTypeError AlertType = "ERROR"

	// AlertTypeWarning is for warning alerts.
	AlertTypeWarning AlertType = "WARNING"

	// AlertTypeInfo is for info alerts.
	AlertTypeInfo AlertType = "INFO"

	// AlertTypeSuccess is for success alerts.
	AlertTypeSuccess AlertType = "SUCCESS"
)

// Alert is an alert message.
type Alert struct {
	// Type is the type of alert.
	Type AlertType `json:"type"`
	// Message is the message for the alert.
	Message string `json:"message"`
	// Expiration is the time the alert expires.
	Expiration *time.Time `json:"expiration,omitempty"`
}

// CreateAlert creates an alert with optional expiration. If the expireAt is < 1
// Expiration will be nil.
func CreateAlert(alertType AlertType, message string, expireAt time.Duration) Alert {
	alert := Alert{
		Type:    alertType,
		Message: message,
	}

	if expireAt > 0 {
		t := time.Now().Add(expireAt)
		alert.Expiration = &t
	}

	return alert
}

type Alerter interface {
	SendAlert(alert Alert)
}

// DispatcherFunc is a function that will be dispatched to handle a payload.
type DispatcherFunc func(ctx context.Context, alerter Alerter, payload Payload) error

// Dispatcher handles actions.
type Dispatcher interface {
	ActionName() string
	Handle(ctx context.Context, alerter Alerter, payload Payload) error
}

// Dispatchers is a slice of Dispatcher.
type Dispatchers []Dispatcher

// ToActionPaths converts Dispatchers to a map.
func (d Dispatchers) ToActionPaths() map[string]DispatcherFunc {
	m := make(map[string]DispatcherFunc)

	for i := range d {
		m[d[i].ActionName()] = d[i].Handle
	}

	return m
}

// Manager manages actions.
type Manager struct {
	logger log.Logger

	// key: string, value: []DispatcherFunc
	dispatches sync.Map
}

// NewManager creates an instance of Manager.
func NewManager(logger log.Logger) *Manager {
	return &Manager{
		logger:     logger.With("component", "action-manager"),
		dispatches: sync.Map{},
	}
}

// Register registers a dispatcher function to an action path.
func (m *Manager) Register(actionPath string, actionFunc DispatcherFunc) error {
	var af []DispatcherFunc

	val, ok := m.dispatches.Load(actionPath)
	if !ok {
		af = make([]DispatcherFunc, 1)
		af[0] = actionFunc
	} else {
		af, ok = val.([]DispatcherFunc)
		if !ok {
			return fmt.Errorf("failed to convert value to []DispatcherFunc")
		}
		af = append(af, actionFunc)
	}

	m.dispatches.Store(actionPath, af)

	return nil
}

// Dispatch dispatches a payload to a path.
func (m *Manager) Dispatch(ctx context.Context, alerter Alerter, actionPath string, payload Payload) error {
	val, ok := m.dispatches.Load(actionPath)
	if !ok {
		return &NotFoundError{Path: actionPath}

	}

	actionFuncs := val.([]DispatcherFunc)
	for _, f := range actionFuncs {
		if err := f(ctx, alerter, payload); err != nil {
			m.logger.Errorf("actionFunc returned err: %s", err)
		}
	}

	return nil
}
