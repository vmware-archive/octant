package astilectron

import (
	"context"

	"github.com/asticode/go-astitools/context"
)

// Session event names
const (
	EventNameSessionCmdClearCache     = "session.cmd.clear.cache"
	EventNameSessionEventClearedCache = "session.event.cleared.cache"
	EventNameSessionEventWillDownload = "session.event.will.download"
)

// Session represents a session
// TODO Add missing session methods
// TODO Add missing session events
// https://github.com/electron/electron/blob/v1.8.1/docs/api/session.md
type Session struct {
	*object
}

// newSession creates a new session
func newSession(parentCtx context.Context, c *asticontext.Canceller, d *dispatcher, i *identifier, w *writer) *Session {
	return &Session{object: newObject(parentCtx, c, d, i, w, i.new())}
}

// ClearCache clears the Session's HTTP cache
func (s *Session) ClearCache() (err error) {
	if err = s.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(s.c, s, s.w, Event{Name: EventNameSessionCmdClearCache, TargetID: s.id}, EventNameSessionEventClearedCache)
	return
}
