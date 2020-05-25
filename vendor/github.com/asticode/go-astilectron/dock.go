package astilectron

import (
	"context"

	"github.com/asticode/go-astikit"
)

// Dock event names
const (
	eventNameDockCmdBounce              = "dock.cmd.bounce"
	eventNameDockCmdBounceDownloads     = "dock.cmd.bounce.downloads"
	eventNameDockCmdCancelBounce        = "dock.cmd.cancel.bounce"
	eventNameDockCmdHide                = "dock.cmd.hide"
	eventNameDockCmdSetBadge            = "dock.cmd.set.badge"
	eventNameDockCmdSetIcon             = "dock.cmd.set.icon"
	eventNameDockCmdShow                = "dock.cmd.show"
	eventNameDockEventBadgeSet          = "dock.event.badge.set"
	eventNameDockEventBouncing          = "dock.event.bouncing"
	eventNameDockEventBouncingCancelled = "dock.event.bouncing.cancelled"
	eventNameDockEventDownloadsBouncing = "dock.event.download.bouncing"
	eventNameDockEventHidden            = "dock.event.hidden"
	eventNameDockEventIconSet           = "dock.event.icon.set"
	eventNameDockEventShown             = "dock.event.shown"
)

// Dock bounce types
const (
	DockBounceTypeCritical      = "critical"
	DockBounceTypeInformational = "informational"
)

// Dock represents a dock
// https://github.com/electron/electron/blob/v1.8.1/docs/api/app.md#appdockbouncetype-macos
type Dock struct {
	*object
}

func newDock(ctx context.Context, d *dispatcher, i *identifier, wrt *writer) *Dock {
	return &Dock{object: newObject(ctx, d, i, wrt, targetIDDock)}
}

// Bounce bounces the dock
func (d *Dock) Bounce(bounceType string) (id int, err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	var e Event
	if e, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdBounce, TargetID: d.id, BounceType: bounceType}, eventNameDockEventBouncing); err != nil {
		return
	}
	if e.ID != nil {
		id = *e.ID
	}
	return
}

// BounceDownloads bounces the downloads part of the dock
func (d *Dock) BounceDownloads(filePath string) (err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdBounceDownloads, TargetID: d.id, FilePath: filePath}, eventNameDockEventDownloadsBouncing)
	return
}

// CancelBounce cancels the dock bounce
func (d *Dock) CancelBounce(id int) (err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdCancelBounce, TargetID: d.id, ID: astikit.IntPtr(id)}, eventNameDockEventBouncingCancelled)
	return
}

// Hide hides the dock
func (d *Dock) Hide() (err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdHide, TargetID: d.id}, eventNameDockEventHidden)
	return
}

// NewMenu creates a new dock menu
func (d *Dock) NewMenu(i []*MenuItemOptions) *Menu {
	return newMenu(d.ctx, d.id, i, d.d, d.i, d.w)
}

// SetBadge sets the badge of the dock
func (d *Dock) SetBadge(badge string) (err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdSetBadge, TargetID: d.id, Badge: badge}, eventNameDockEventBadgeSet)
	return
}

// SetIcon sets the icon of the dock
func (d *Dock) SetIcon(image string) (err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdSetIcon, TargetID: d.id, Image: image}, eventNameDockEventIconSet)
	return
}

// Show shows the dock
func (d *Dock) Show() (err error) {
	if err = d.ctx.Err(); err != nil {
		return
	}
	_, err = synchronousEvent(d.ctx, d, d.w, Event{Name: eventNameDockCmdShow, TargetID: d.id}, eventNameDockEventShown)
	return
}
