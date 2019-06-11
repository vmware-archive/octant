package componentcache

import (
	"context"
	"fmt"
	"testing"

	"github.com/heptio/developer-dash/pkg/view/component"

	"github.com/stretchr/testify/assert"

	lru "github.com/hashicorp/golang-lru"
)

func Test_ComponentCache_Update(t *testing.T) {
	components, err := lru.New(10)
	assert.NoError(t, err)

	ctx := context.Background()
	cc := &componentCache{
		components: components,
		ch:         make(chan Event, 1),
	}

	scenarios := []struct {
		name     string
		updateFn UpdateFn
		contains string
		err      error
	}{
		{
			name:     "success",
			contains: "success",
			err:      nil,
			updateFn: func(ctx context.Context, ch chan Event) (string, error) {
				key := "success"
				comp := component.NewText("success")

				cc.Add(key, comp)
				return key, nil
			},
		},
		{
			name:     "error (not found)",
			contains: "not found",
			err:      fmt.Errorf("%s not found in ComponentCache", "error"),
			updateFn: func(ctx context.Context, ch chan Event) (string, error) {
				key := "error"
				return key, nil
			},
		},
		{
			name:     "error (update error)",
			contains: "bad update",
			err:      fmt.Errorf("bad update %s", "bad_key"),
			updateFn: func(ctx context.Context, ch chan Event) (string, error) {
				key := "bad_key"
				return key, fmt.Errorf("bad update %s", key)
			},
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			comp, err := cc.Update(ctx, ts.updateFn)
			assert.Equal(t, ts.err, err)
			if ts.err == nil {
				assert.Contains(t, ts.contains, comp.String())
			}
		})
	}
}

func Test_ComponentCache_Get_Add(t *testing.T) {
	components, err := lru.New(1)
	assert.NoError(t, err)

	cc := &componentCache{
		components: components,
		ch:         make(chan Event, 1),
	}

	key1 := "testKey1"
	key2 := "testKey2"
	comp := component.NewText("testText")

	// empty cache
	value, ok := cc.Get(key1)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, value)

	// key1, no eviction
	eviction := cc.Add(key1, comp)
	assert.Equal(t, false, eviction)

	// key2, eviction
	eviction = cc.Add(key2, comp)
	assert.Equal(t, true, eviction)

	// key1 was evicted
	value, ok = cc.Get(key1)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, value)

	// key2 value
	value, ok = cc.Get(key2)
	assert.Equal(t, true, ok)
	assert.Equal(t, comp, value)
}

func Test_Worker_Event(t *testing.T) {
	components, err := lru.New(1)
	assert.NoError(t, err)

	cc := &componentCache{
		components: components,
		ch:         make(chan Event, 1),
	}

	scenarios := []struct {
		name     string
		expected bool
		ctx      context.Context

		eventKey   string
		eventValue component.Component
		eventError error
	}{
		{
			name:       "<- event (no error)",
			expected:   true,
			ctx:        context.Background(),
			eventKey:   "successKey",
			eventValue: component.NewText("testSuccess"),
			eventError: nil,
		},
		{
			name:       "<- event (error)",
			expected:   true,
			ctx:        context.Background(),
			eventKey:   "errorKey",
			eventValue: component.NewText("testError"),
			eventError: fmt.Errorf("%s", "error event"),
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			e := Event{
				Name:      ts.name,
				Key:       ts.eventKey,
				Component: ts.eventValue,
				Err:       ts.eventError,
			}
			cc.ch <- e

			value := worker(ts.ctx, cc, cc.ch)
			assert.Equal(t, ts.expected, value)

			comp, ok := cc.Get(ts.eventKey)
			assert.Equal(t, ok, true)

			if ts.eventError != nil {
				assert.Contains(t, comp.String(), ts.eventError.Error())
			} else {
				assert.Equal(t, ts.eventValue, comp)
			}
		})
	}
}
func Test_Worker_Done_Skip(t *testing.T) {
	components, err := lru.New(1)
	assert.NoError(t, err)

	cc := &componentCache{
		components: components,
		ch:         make(chan Event, 1),
	}

	scenarios := []struct {
		name     string
		expected bool
		ctx      context.Context
	}{
		{
			name:     "ctx.Done",
			expected: false,
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
		},
		{
			name:     "default",
			expected: true,
			ctx:      context.Background(),
		},
	}

	for _, ts := range scenarios {
		t.Run(ts.name, func(t *testing.T) {
			value := worker(ts.ctx, cc, cc.ch)
			assert.Equal(t, ts.expected, value)
		})
	}
}
