package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/heptio/developer-dash/internal/cluster"
	"github.com/heptio/developer-dash/internal/clustereye"
)

type namespacesResponse struct {
	Namespaces []string `json:"namespaces,omitempty"`
}

// NamespacesGenerator generates namespaces events.
type NamespacesGenerator struct {
	// NamespaceClient is a namespaces client.
	NamespaceClient cluster.NamespaceInterface
}

var _ clustereye.Generator = (*NamespacesGenerator)(nil)

// Event generates namespaces events
func (g *NamespacesGenerator) Event(ctx context.Context) (clustereye.Event, error) {
	if g.NamespaceClient == nil {
		return clustereye.Event{}, errors.New("unable to query namespaces, client is nil")
	}

	names, err := g.NamespaceClient.Names()
	if err != nil {
		initialNamespace := g.NamespaceClient.InitialNamespace()
		names = []string{initialNamespace}
	}

	nr := &namespacesResponse{Namespaces: names}
	data, err := json.Marshal(nr)
	if err != nil {
		return clustereye.Event{}, errors.New("unable to marshal namespaces")
	}

	return clustereye.Event{
		Type: clustereye.EventTypeNamespaces,
		Data: data,
	}, nil
}

// ScheduleDelay returns how long to delay before running this generator again.
func (NamespacesGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

// Name returns the generator's name.
func (NamespacesGenerator) Name() string {
	return "namespaces"
}
