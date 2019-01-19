package astilectron

import (
	"context"

	"github.com/asticode/go-astitools/context"
)

// Menu event names
const (
	EventNameMenuCmdCreate      = "menu.cmd.create"
	EventNameMenuCmdDestroy     = "menu.cmd.destroy"
	EventNameMenuEventCreated   = "menu.event.created"
	EventNameMenuEventDestroyed = "menu.event.destroyed"
)

// Menu represents a menu
// https://github.com/electron/electron/blob/v1.8.1/docs/api/menu.md
type Menu struct {
	*subMenu
}

// newMenu creates a new menu
func newMenu(parentCtx context.Context, rootID string, items []*MenuItemOptions, c *asticontext.Canceller, d *dispatcher, i *identifier, w *writer) (m *Menu) {
	// Init
	m = &Menu{newSubMenu(parentCtx, rootID, items, c, d, i, w)}

	// Make sure the menu's context is cancelled once the destroyed event is received
	m.On(EventNameMenuEventDestroyed, func(e Event) (deleteListener bool) {
		m.cancel()
		return true
	})
	return
}

// toEvent returns the menu in the proper event format
func (m *Menu) toEvent() *EventMenu {
	return &EventMenu{m.subMenu.toEvent()}
}

// Create creates the menu
func (m *Menu) Create() (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(m.c, m, m.w, Event{Name: EventNameMenuCmdCreate, TargetID: m.id, Menu: m.toEvent()}, EventNameMenuEventCreated)
	return
}

// Destroy destroys the menu
func (m *Menu) Destroy() (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(m.c, m, m.w, Event{Name: EventNameMenuCmdDestroy, TargetID: m.id, Menu: m.toEvent()}, EventNameMenuEventDestroyed)
	return
}
