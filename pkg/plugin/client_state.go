package plugin

import (
	"context"

	ocontext "github.com/vmware-tanzu/octant/internal/context"
)

type ClientState interface {
	ClientID() string
	Filters() []Filter
	Namespace() string
	ContextName() string
}

func ClientStateFrom(ctx context.Context) clientState {
	cs := ocontext.ClientStateFrom(ctx)
	filters := []Filter{}

	for _, f := range cs.Filters {
		filters = append(filters, clientFilter{key: f.Key, value: f.Value})
	}

	return clientState{
		clientID:    cs.ClientID,
		filters:     filters,
		namespace:   cs.Namespace,
		contextName: cs.ContextName,
	}
}

type Filter interface {
	Key() string
	Value() string
}

var _ ClientState = (*clientState)(nil)

type clientState struct {
	clientID    string
	filters     []Filter
	namespace   string
	contextName string
}

func (c clientState) ClientID() string {
	return c.clientID
}
func (c clientState) Filters() []Filter {
	return c.filters
}
func (c clientState) Namespace() string {
	return c.namespace
}
func (c clientState) ContextName() string {
	return c.contextName
}

var _ Filter = (*clientFilter)(nil)

type clientFilter struct {
	key   string
	value string
}

func (cf clientFilter) Key() string {
	return cf.key
}

func (cf clientFilter) Value() string {
	return cf.value
}
