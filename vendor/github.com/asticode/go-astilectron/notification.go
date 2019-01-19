package astilectron

import "github.com/asticode/go-astitools/context"

// Notification event names
const (
	eventNameNotificationCmdCreate    = "notification.cmd.create"
	eventNameNotificationCmdShow      = "notification.cmd.show"
	EventNameNotificationEventClicked = "notification.event.clicked"
	EventNameNotificationEventClosed  = "notification.event.closed"
	EventNameNotificationEventCreated = "notification.event.created"
	EventNameNotificationEventReplied = "notification.event.replied"
	EventNameNotificationEventShown   = "notification.event.shown"
)

// Notification represents a notification
// https://github.com/electron/electron/blob/v1.8.1/docs/api/notification.md
type Notification struct {
	isSupported bool
	o           *NotificationOptions
	*object
}

// NotificationOptions represents notification options
type NotificationOptions struct {
	Body             string `json:"body,omitempty"`
	HasReply         *bool  `json:"hasReply,omitempty"`
	Icon             string `json:"icon,omitempty"`
	ReplyPlaceholder string `json:"replyPlaceholder,omitempty"`
	Silent           *bool  `json:"silent,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Subtitle         string `json:"subtitle,omitempty"`
	Title            string `json:"title,omitempty"`
}

func newNotification(o *NotificationOptions, isSupported bool, c *asticontext.Canceller, d *dispatcher, i *identifier, wrt *writer) *Notification {
	return &Notification{
		isSupported: isSupported,
		o:           o,
		object:      newObject(nil, c, d, i, wrt, i.new()),
	}
}

// Create creates the notification
func (n *Notification) Create() (err error) {
	if !n.isSupported {
		return
	}
	if err = n.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(n.c, n, n.w, Event{Name: eventNameNotificationCmdCreate, TargetID: n.id, NotificationOptions: n.o}, EventNameNotificationEventCreated)
	return
}

// Show shows the notification
func (n *Notification) Show() (err error) {
	if !n.isSupported {
		return
	}
	if err = n.isActionable(); err != nil {
		return
	}
	_, err = synchronousEvent(n.c, n, n.w, Event{Name: eventNameNotificationCmdShow, TargetID: n.id}, EventNameNotificationEventShown)
	return
}
