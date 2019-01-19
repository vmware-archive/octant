package astilectron

import (
	"context"

	"github.com/asticode/go-astitools/context"
)

// Menu item event names
const (
	EventNameMenuItemCmdSetChecked   = "menu.item.cmd.set.checked"
	EventNameMenuItemCmdSetEnabled   = "menu.item.cmd.set.enabled"
	EventNameMenuItemCmdSetLabel     = "menu.item.cmd.set.label"
	EventNameMenuItemCmdSetVisible   = "menu.item.cmd.set.visible"
	EventNameMenuItemEventCheckedSet = "menu.item.event.checked.set"
	EventNameMenuItemEventClicked    = "menu.item.event.clicked"
	EventNameMenuItemEventEnabledSet = "menu.item.event.enabled.set"
	EventNameMenuItemEventLabelSet   = "menu.item.event.label.set"
	EventNameMenuItemEventVisibleSet = "menu.item.event.visible.set"
)

// Menu item roles
var (
	// All
	MenuItemRoleClose              = PtrStr("close")
	MenuItemRoleCopy               = PtrStr("copy")
	MenuItemRoleCut                = PtrStr("cut")
	MenuItemRoleDelete             = PtrStr("delete")
	MenuItemRoleEditMenu           = PtrStr("editMenu")
	MenuItemRoleForceReload        = PtrStr("forcereload")
	MenuItemRoleMinimize           = PtrStr("minimize")
	MenuItemRolePaste              = PtrStr("paste")
	MenuItemRolePasteAndMatchStyle = PtrStr("pasteandmatchstyle")
	MenuItemRoleQuit               = PtrStr("quit")
	MenuItemRoleRedo               = PtrStr("redo")
	MenuItemRoleReload             = PtrStr("reload")
	MenuItemRoleResetZoom          = PtrStr("resetzoom")
	MenuItemRoleSelectAll          = PtrStr("selectall")
	MenuItemRoleToggleDevTools     = PtrStr("toggledevtools")
	MenuItemRoleToggleFullScreen   = PtrStr("togglefullscreen")
	MenuItemRoleUndo               = PtrStr("undo")
	MenuItemRoleWindowMenu         = PtrStr("windowMenu")
	MenuItemRoleZoomOut            = PtrStr("zoomout")
	MenuItemRoleZoomIn             = PtrStr("zoomin")

	// MacOSX
	MenuItemRoleAbout         = PtrStr("about")
	MenuItemRoleHide          = PtrStr("hide")
	MenuItemRoleHideOthers    = PtrStr("hideothers")
	MenuItemRoleUnhide        = PtrStr("unhide")
	MenuItemRoleStartSpeaking = PtrStr("startspeaking")
	MenuItemRoleStopSpeaking  = PtrStr("stopspeaking")
	MenuItemRoleFront         = PtrStr("front")
	MenuItemRoleZoom          = PtrStr("zoom")
	MenuItemRoleWindow        = PtrStr("window")
	MenuItemRoleHelp          = PtrStr("help")
	MenuItemRoleServices      = PtrStr("services")
)

// Menu item types
var (
	MenuItemTypeNormal    = PtrStr("normal")
	MenuItemTypeSeparator = PtrStr("separator")
	MenuItemTypeCheckbox  = PtrStr("checkbox")
	MenuItemTypeRadio     = PtrStr("radio")
)

// MenuItem represents a menu item
type MenuItem struct {
	*object
	o *MenuItemOptions
	// We must store the root ID since everytime we update a sub menu we need to set the root menu all over again in electron
	rootID string
	s      *SubMenu
}

// MenuItemOptions represents menu item options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use PtrBool, PtrInt or PtrStr
// to fill the struct
// https://github.com/electron/electron/blob/v1.8.1/docs/api/menu-item.md
type MenuItemOptions struct {
	Accelerator *Accelerator       `json:"accelerator,omitempty"`
	Checked     *bool              `json:"checked,omitempty"`
	Enabled     *bool              `json:"enabled,omitempty"`
	Icon        *string            `json:"icon,omitempty"`
	Label       *string            `json:"label,omitempty"`
	OnClick     Listener           `json:"-"`
	Position    *string            `json:"position,omitempty"`
	Role        *string            `json:"role,omitempty"`
	SubLabel    *string            `json:"sublabel,omitempty"`
	SubMenu     []*MenuItemOptions `json:"-"`
	Type        *string            `json:"type,omitempty"`
	Visible     *bool              `json:"visible,omitempty"`
}

// newMenu creates a new menu item
func newMenuItem(parentCtx context.Context, rootID string, o *MenuItemOptions, c *asticontext.Canceller, d *dispatcher, i *identifier, w *writer) (m *MenuItem) {
	m = &MenuItem{
		o:      o,
		object: newObject(parentCtx, c, d, i, w, i.new()),
		rootID: rootID,
	}
	if o.OnClick != nil {
		m.On(EventNameMenuItemEventClicked, o.OnClick)
	}
	if len(o.SubMenu) > 0 {
		m.s = &SubMenu{newSubMenu(parentCtx, rootID, o.SubMenu, c, d, i, w)}
	}
	return
}

// toEvent returns the menu item in the proper event format
func (i *MenuItem) toEvent() (e *EventMenuItem) {
	e = &EventMenuItem{
		ID:      i.id,
		Options: i.o,
		RootID:  i.rootID,
	}
	if i.s != nil {
		e.SubMenu = i.s.toEvent()
	}
	return
}

// SubMenu returns the menu item sub menu
func (i *MenuItem) SubMenu() *SubMenu {
	return i.s
}

// SetChecked sets the checked attribute
func (i *MenuItem) SetChecked(checked bool) (err error) {
	if err = i.isActionable(); err != nil {
		return
	}
	i.o.Checked = PtrBool(checked)
	_, err = synchronousEvent(i.c, i, i.w, Event{Name: EventNameMenuItemCmdSetChecked, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Checked: i.o.Checked}}, EventNameMenuItemEventCheckedSet)
	return
}

// SetEnabled sets the enabled attribute
func (i *MenuItem) SetEnabled(enabled bool) (err error) {
	if err = i.isActionable(); err != nil {
		return
	}
	i.o.Enabled = PtrBool(enabled)
	_, err = synchronousEvent(i.c, i, i.w, Event{Name: EventNameMenuItemCmdSetEnabled, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Enabled: i.o.Enabled}}, EventNameMenuItemEventEnabledSet)
	return
}

// SetLabel sets the label attribute
func (i *MenuItem) SetLabel(label string) (err error) {
	if err = i.isActionable(); err != nil {
		return
	}
	i.o.Label = PtrStr(label)
	_, err = synchronousEvent(i.c, i, i.w, Event{Name: EventNameMenuItemCmdSetLabel, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Label: i.o.Label}}, EventNameMenuItemEventLabelSet)
	return
}

// SetVisible sets the visible attribute
func (i *MenuItem) SetVisible(visible bool) (err error) {
	if err = i.isActionable(); err != nil {
		return
	}
	i.o.Visible = PtrBool(visible)
	_, err = synchronousEvent(i.c, i, i.w, Event{Name: EventNameMenuItemCmdSetVisible, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Visible: i.o.Visible}}, EventNameMenuItemEventVisibleSet)
	return
}
