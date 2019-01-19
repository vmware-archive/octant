package astilectron

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/asticode/go-astitools/context"
)

// Sub menu event names
const (
	EventNameSubMenuCmdAppend        = "sub.menu.cmd.append"
	EventNameSubMenuCmdClosePopup    = "sub.menu.cmd.close.popup"
	EventNameSubMenuCmdInsert        = "sub.menu.cmd.insert"
	EventNameSubMenuCmdPopup         = "sub.menu.cmd.popup"
	EventNameSubMenuEventAppended    = "sub.menu.event.appended"
	EventNameSubMenuEventClosedPopup = "sub.menu.event.closed.popup"
	EventNameSubMenuEventInserted    = "sub.menu.event.inserted"
	EventNameSubMenuEventPoppedUp    = "sub.menu.event.popped.up"
)

// SubMenu represents an exported sub menu
type SubMenu struct {
	*subMenu
}

// subMenu represents an internal sub menu
// We use this internal subMenu in SubMenu and Menu since all functions of subMenu should be in SubMenu and Menu but
// some functions of Menu shouldn't be in SubMenu and vice versa
type subMenu struct {
	*object
	items []*MenuItem
	// We must store the root ID since everytime we update a sub menu we need to set the root menu all over again in electron
	rootID string
}

// newSubMenu creates a new sub menu
func newSubMenu(parentCtx context.Context, rootID string, items []*MenuItemOptions, c *asticontext.Canceller, d *dispatcher, i *identifier, w *writer) *subMenu {
	// Init
	var m = &subMenu{
		object: newObject(parentCtx, c, d, i, w, i.new()),
		rootID: rootID,
	}

	// Parse items
	for _, o := range items {
		m.items = append(m.items, newMenuItem(m.ctx, rootID, o, c, d, i, w))
	}
	return m
}

// toEvent returns the sub menu in the proper event format
func (m *subMenu) toEvent() (e *EventSubMenu) {
	e = &EventSubMenu{
		ID:     m.id,
		RootID: m.rootID,
	}
	for _, i := range m.items {
		e.Items = append(e.Items, i.toEvent())
	}
	return
}

// NewItem returns a new menu item
func (m *subMenu) NewItem(o *MenuItemOptions) *MenuItem {
	return newMenuItem(m.ctx, m.rootID, o, m.c, m.d, m.i, m.w)
}

// SubMenu returns the sub menu at the specified indexes
func (m *subMenu) SubMenu(indexes ...int) (s *SubMenu, err error) {
	var is = m
	var processedIndexes = []string{}
	for _, index := range indexes {
		if index >= len(is.items) {
			return nil, fmt.Errorf("Submenu at %s has %d items, invalid index %d", strings.Join(processedIndexes, ":"), len(is.items), index)
		}
		s = is.items[index].s
		processedIndexes = append(processedIndexes, strconv.Itoa(index))
		if s == nil {
			return nil, fmt.Errorf("No submenu at %s", strings.Join(processedIndexes, ":"))
		}
		is = s.subMenu
	}
	return
}

// Item returns the item at the specified indexes
func (m *subMenu) Item(indexes ...int) (mi *MenuItem, err error) {
	var is = m
	if len(indexes) > 1 {
		var s *SubMenu
		if s, err = m.SubMenu(indexes[:len(indexes)-1]...); err != nil {
			return
		}
		is = s.subMenu
	}
	var index = indexes[len(indexes)-1]
	if index >= len(is.items) {
		return nil, fmt.Errorf("Submenu has %d items, invalid index %d", len(is.items), index)
	}
	mi = is.items[index]
	return
}

// Append appends a menu item into the sub menu
func (m *subMenu) Append(i *MenuItem) (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	if _, err = synchronousEvent(m.c, m, m.w, Event{Name: EventNameSubMenuCmdAppend, TargetID: m.id, MenuItem: i.toEvent()}, EventNameSubMenuEventAppended); err != nil {
		return
	}
	m.items = append(m.items, i)
	return
}

// Insert inserts a menu item to the position of the sub menu
func (m *subMenu) Insert(pos int, i *MenuItem) (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	if pos > len(m.items) {
		err = fmt.Errorf("Submenu has %d items, position %d is invalid", len(m.items), pos)
		return
	}
	if _, err = synchronousEvent(m.c, m, m.w, Event{Name: EventNameSubMenuCmdInsert, TargetID: m.id, MenuItem: i.toEvent(), MenuItemPosition: PtrInt(pos)}, EventNameSubMenuEventInserted); err != nil {
		return
	}
	m.items = append(m.items[:pos], append([]*MenuItem{i}, m.items[pos:]...)...)
	return
}

// MenuPopupOptions represents menu pop options
type MenuPopupOptions struct {
	PositionOptions
	PositioningItem *int `json:"positioningItem,omitempty"`
}

// Popup pops up the menu as a context menu in the focused window
func (m *subMenu) Popup(o *MenuPopupOptions) error {
	return m.PopupInWindow(nil, o)
}

// PopupInWindow pops up the menu as a context menu in the specified window
func (m *subMenu) PopupInWindow(w *Window, o *MenuPopupOptions) (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	var e = Event{Name: EventNameSubMenuCmdPopup, TargetID: m.id, MenuPopupOptions: o}
	if w != nil {
		e.WindowID = w.id
	}
	_, err = synchronousEvent(m.c, m, m.w, e, EventNameSubMenuEventPoppedUp)
	return
}

// ClosePopup close the context menu in the focused window
func (m *subMenu) ClosePopup() error {
	return m.ClosePopupInWindow(nil)
}

// ClosePopupInWindow close the context menu in the specified window
func (m *subMenu) ClosePopupInWindow(w *Window) (err error) {
	if err = m.isActionable(); err != nil {
		return
	}
	var e = Event{Name: EventNameSubMenuCmdClosePopup, TargetID: m.id}
	if w != nil {
		e.WindowID = w.id
	}
	_, err = synchronousEvent(m.c, m, m.w, e, EventNameSubMenuEventClosedPopup)
	return
}
