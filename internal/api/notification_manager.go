package api

import (
	"context"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/api"
)

const (
	RequestNotifications = "event.octant.dev/notification"
)

type notificationStateManager struct {
	client api.OctantClient
	config config.Dash
	ctx    context.Context

	lastMessageSend int
}

var _ StateManager = (*notificationStateManager)(nil)

// NewNotificationsStateManager returns a terminal state manager.
func NewNotificationsStateManager(dashConfig config.Dash) *notificationStateManager {
	return &notificationStateManager{
		config:          dashConfig,
		lastMessageSend: 0,
	}
}

func (s *notificationStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestNotifications,
			Handler:     s.getNotifications,
		},
	}
}

func (s *notificationStateManager) getNotifications(_ octant.State, payload action.Payload) error {
	length := len(s.config.ErrorStore().List())

	if s.lastMessageSend <= length {
		for ; s.lastMessageSend <= length; s.lastMessageSend++ {
			newEvent := event.Event{
				Type: event.EventTypeNotification,
				Data: map[string]interface{}{
					"error": s.config.ErrorStore().List()[s.lastMessageSend].Error(),
					"name":  s.config.ErrorStore().List()[s.lastMessageSend].Name(),
				},
			}
			s.client.Send(newEvent)
		}
	}

	return nil
}

func (s *notificationStateManager) Start(ctx context.Context, _ octant.State, client api.OctantClient) {
	s.client = client
	s.ctx = ctx
}
