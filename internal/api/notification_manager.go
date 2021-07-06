package api

import (
	"context"

	"github.com/vmware-tanzu/octant/internal/errors"

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
}

var _ StateManager = (*notificationStateManager)(nil)
var _ errors.Observer = (*notificationStateManager)(nil)

// NewNotificationsStateManager returns a terminal state manager.
func NewNotificationsStateManager(dashConfig config.Dash) *notificationStateManager {

	nsm := notificationStateManager{
		config: dashConfig,
	}
	dashConfig.ErrorStore().Subscribe(&nsm)
	return &nsm
}

func (s *notificationStateManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: RequestNotifications,
			Handler:     s.getNotifications,
		},
	}
}

// Always send the entire list to the user because the back-end is responsible
// to take care of duplicated event and because ErrorStore was a max size of 50
func (s *notificationStateManager) getNotifications(_ octant.State, _ action.Payload) error {
	list := s.Marshal(s.config.ErrorStore().List())

	newEvent := event.Event{
		Type: event.EventTypeNotification,
		Data: map[string]interface{}{
			"errors": list,
		},
	}
	s.client.Send(newEvent)

	return nil
}

func (s *notificationStateManager) Marshal(ie []errors.InternalError) []map[string]interface{} {
	var result []map[string]interface{}
	size := len(ie)

	for i := 0; i < size; i++ {
		result = append(result, map[string]interface{}{
			"error":     s.config.ErrorStore().List()[i].Error(),
			"name":      s.config.ErrorStore().List()[i].Name(),
			"timestamp": s.config.ErrorStore().List()[i].Timestamp(),
		})
	}

	return result
}

func (s *notificationStateManager) Update() {
	s.getNotifications(nil, nil)
}

func (s *notificationStateManager) Start(ctx context.Context, _ octant.State, client api.OctantClient) {
	s.client = client
	s.ctx = ctx
}
