package componentcache

import (
	"context"
	"fmt"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/heptio/developer-dash/pkg/view/component"
)

//go:generate mockgen -source=componentcache.go -destination=./fake/mock_component_cache.go -package=fake github.com/heptio/developer-dash/internal/componentcache/ ComponentCache

// Event holds a Key, Component, and Error and is for writing to the cache.
type Event struct {
	Key        string
	Name       string
	CComponent component.Component
	Err        error
}

// UpdateFn provides the current context and event channel and should return the key for the component.
type UpdateFn func(context.Context, chan Event) (string, error)

// ComponentCache is cache of Components
type ComponentCache interface {
	Add(key, value interface{}) bool
	Get(key interface{}) (component.Component, bool)
	Update(context.Context, UpdateFn) (component.Component, error)
	Start(context.Context)
}

type componentCache struct {
	components *lru.Cache
	ch         chan Event
}

// NewComponentCache creates a new component cache.
func NewComponentCache(ctx context.Context) (ComponentCache, error) {
	components, err := lru.New(500)
	if err != nil {
		return nil, err
	}

	cc := &componentCache{
		components: components,
		ch:         make(chan Event, 1),
	}
	cc.Start(ctx)
	return cc, nil
}

// Start starts the event loop on the ComponentCache
func (cc *componentCache) Start(ctx context.Context) {
	shouldRun := true
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			if !shouldRun {
				break
			}
			shouldRun = worker(ctx, cc, cc.ch)
		}
	}()
}

func worker(ctx context.Context, cc *componentCache, ch chan Event) bool {
	select {
	case <-ctx.Done():
		return false
	case e := <-ch:
		if e.Err != nil {
			title := component.Title(component.NewText(e.Name))
			errComponent := component.NewError(title, e.Err)
			cc.Add(e.Key, errComponent)
			return true
		}
		cc.Add(e.Key, e.CComponent)
		return true
	default:
		return true
	}
}

// Add inserts a value in to the ComponentCache
func (cc *componentCache) Add(key interface{}, value interface{}) bool {
	return cc.components.Add(key, value)
}

// Get fetches a value from  the ComponentCache
func (cc *componentCache) Get(key interface{}) (component.Component, bool) {
	v, ok := cc.components.Get(key)
	if !ok {
		return nil, ok
	}
	return v.(component.Component), ok
}

// Update launches the updateFn and returns the current value from the ComponentCache
func (cc *componentCache) Update(ctx context.Context, updateFn UpdateFn) (component.Component, error) {
	key, err := updateFn(ctx, cc.ch)
	if err != nil {
		return nil, err
	}

	ccomponent, ok := cc.Get(key)
	if !ok {
		return ccomponent, fmt.Errorf("%s not found in ComponentCache", key)
	}
	return ccomponent, nil
}
