package astilectron

import (
	"context"

	"github.com/asticode/go-astikit"
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
	MenuItemRoleClose              = astikit.StrPtr("close")
	MenuItemRoleCopy               = astikit.StrPtr("copy")
	MenuItemRoleCut                = astikit.StrPtr("cut")
	MenuItemRoleDelete             = astikit.StrPtr("delete")
	MenuItemRoleEditMenu           = astikit.StrPtr("editMenu")
	MenuItemRoleForceReload        = astikit.StrPtr("forcereload")
	MenuItemRoleMinimize           = astikit.StrPtr("minimize")
	MenuItemRolePaste              = astikit.StrPtr("paste")
	MenuItemRolePasteAndMatchStyle = astikit.StrPtr("pasteandmatchstyle")
	MenuItemRoleQuit               = astikit.StrPtr("quit")
	MenuItemRoleRedo               = astikit.StrPtr("redo")
	MenuItemRoleReload             = astikit.StrPtr("reload")
	MenuItemRoleResetZoom          = astikit.StrPtr("resetzoom")
	MenuItemRoleSelectAll          = astikit.StrPtr("selectall")
	MenuItemRoleToggleDevTools     = astikit.StrPtr("toggledevtools")
	MenuItemRoleToggleFullScreen   = astikit.StrPtr("togglefullscreen")
	MenuItemRoleUndo               = astikit.StrPtr("undo")
	MenuItemRoleWindowMenu         = astikit.StrPtr("windowMenu")
	MenuItemRoleZoomOut            = astikit.StrPtr("zoomout")
	MenuItemRoleZoomIn             = astikit.StrPtr("zoomin")

	// MacOSX
	MenuItemRoleAbout         = astikit.StrPtr("about")
	MenuItemRoleHide          = astikit.StrPtr("hide")
	MenuItemRoleHideOthers    = astikit.StrPtr("hideothers")
	MenuItemRoleUnhide        = astikit.StrPtr("unhide")
	MenuItemRoleStartSpeaking = astikit.StrPtr("startspeaking")
	MenuItemRoleStopSpeaking  = astikit.StrPtr("stopspeaking")
	MenuItemRoleFront         = astikit.StrPtr("front")
	MenuItemRoleZoom          = astikit.StrPtr("zoom")
	MenuItemRoleWindow        = astikit.StrPtr("window")
	MenuItemRoleHelp          = astikit.StrPtr("help")
	MenuItemRoleServices      = astikit.StrPtr("services")
)

// Menu item types
var (
	MenuItemTypeNormal    = astikit.StrPtr("normal")
	MenuItemTypeSeparator = astikit.StrPtr("separator")
	MenuItemTypeCheckbox  = astikit.StrPtr("checkbox")
	MenuItemTypeRadio     = astikit.StrPtr("radio")
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
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use astikit.BoolPtr, astikit.IntPtr or astikit.StrPtr
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
func newMenuItem(ctx context.Context, rootID string, o *MenuItemOptions, d *dispatcher, i *identifier, w *writer) (m *MenuItem) {
	m = &MenuItem{
		o:      o,
		object: newObject(ctx, d, i, w, i.new()),
		rootID: rootID,
	}
	if o.OnClick != nil {
		m.On(EventNameMenuItemEventClicked, o.OnClick)
	}
	if len(o.SubMenu) > 0 {
		m.s = &SubMenu{newSubMenu(ctx, rootID, o.SubMenu, d, i, w)}
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
	if err = i.ctx.Err(); err != nil {
		return
	}
	i.o.Checked = astikit.BoolPtr(checked)
	_, err = synchronousEvent(i.ctx, i, i.w, Event{Name: EventNameMenuItemCmdSetChecked, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Checked: i.o.Checked}}, EventNameMenuItemEventCheckedSet)
	return
}

// SetEnabled sets the enabled attribute
func (i *MenuItem) SetEnabled(enabled bool) (err error) {
	if err = i.ctx.Err(); err != nil {
		return
	}
	i.o.Enabled = astikit.BoolPtr(enabled)
	_, err = synchronousEvent(i.ctx, i, i.w, Event{Name: EventNameMenuItemCmdSetEnabled, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Enabled: i.o.Enabled}}, EventNameMenuItemEventEnabledSet)
	return
}

// SetLabel sets the label attribute
func (i *MenuItem) SetLabel(label string) (err error) {
	if err = i.ctx.Err(); err != nil {
		return
	}
	i.o.Label = astikit.StrPtr(label)
	_, err = synchronousEvent(i.ctx, i, i.w, Event{Name: EventNameMenuItemCmdSetLabel, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Label: i.o.Label}}, EventNameMenuItemEventLabelSet)
	return
}

// SetVisible sets the visible attribute
func (i *MenuItem) SetVisible(visible bool) (err error) {
	if err = i.ctx.Err(); err != nil {
		return
	}
	i.o.Visible = astikit.BoolPtr(visible)
	_, err = synchronousEvent(i.ctx, i, i.w, Event{Name: EventNameMenuItemCmdSetVisible, TargetID: i.id, MenuItemOptions: &MenuItemOptions{Visible: i.o.Visible}}, EventNameMenuItemEventVisibleSet)
	return
}
