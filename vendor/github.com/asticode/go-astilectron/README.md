[![GoReportCard](http://goreportcard.com/badge/github.com/asticode/go-astilectron)](http://goreportcard.com/report/github.com/asticode/go-astilectron)
[![GoDoc](https://godoc.org/github.com/asticode/go-astilectron?status.svg)](https://godoc.org/github.com/asticode/go-astilectron)
[![Travis](https://travis-ci.org/asticode/go-astilectron.svg?branch=master)](https://travis-ci.org/asticode/go-astilectron#)
[![Coveralls](https://coveralls.io/repos/github/asticode/go-astilectron/badge.svg?branch=master)](https://coveralls.io/repos/github/asticode/go-astilectron)

Thanks to `go-astilectron` build cross platform GUI apps with GO and HTML/JS/CSS. It is the official GO bindings of [astilectron](https://github.com/asticode/astilectron) and is powered by [Electron](https://github.com/electron/electron).

# Demo

To see a minimal Astilectron app, checkout out the [demo](https://github.com/asticode/go-astilectron-demo).

It uses the [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) and the [bundler](https://github.com/asticode/go-astilectron-bundler).

# Real-life examples

Here's a list of awesome projects using `go-astilectron` (if you're using `go-astilectron` and want your project to be listed here please submit a PR):

- [go-astivid](https://github.com/asticode/go-astivid) Video tools written in GO
- [GroupMatcher](https://github.com/veecue/GroupMatcher) Program to allocate persons to groups while trying to fulfill all the given wishes as good as possible
- [ipeye-onvif](https://github.com/deepch/ipeye-onvif) ONVIF Search Tool
- [Stellite GUI Miner](https://github.com/stellitecoin/GUI-miner) An easy to use GUI cryptocurrency miner for Stellite


# Bootstrap

For convenience purposes, a [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) has been implemented.

The bootstrap allows you to quickly create a one-window application.

There's no obligation to use it, but it's strongly recommended.

If you decide to use it, read thoroughly the documentation as you'll have to structure your project in a specific way.

# Bundler

Still for convenience purposes, a [bundler](https://github.com/asticode/go-astilectron-bundler) has been implemented.

The bundler allows you to bundle your app for every os/arch combinations and get a nice set of files to send your users.

# Quick start

WARNING: the code below doesn't handle errors for readibility purposes. However you SHOULD!

## Import `go-astilectron`

To import `go-astilectron` run:

    $ go get -u github.com/asticode/go-astilectron

## Start `go-astilectron`

```go
// Initialize astilectron
var a, _ = astilectron.New(astilectron.Options{
    AppName: "<your app name>",
    AppIconDefaultPath: "<your .png icon>", // If path is relative, it must be relative to the data directory
    AppIconDarwinPath:  "<your .icns icon>", // Same here
    BaseDirectoryPath: "<where you want the provisioner to install the dependencies>",
})
defer a.Close()

// Start astilectron
a.Start()

// Blocking pattern
a.Wait()
```

For everything to work properly we need to fetch 2 dependencies : [astilectron](https://github.com/asticode/astilectron) and [Electron](https://github.com/electron/electron). `.Start()` takes care of it by downloading the sources and setting them up properly.

In case you want to embed the sources in the binary to keep a unique binary you can use the **NewDisembedderProvisioner** function to get the proper **Provisioner** and attach it to `go-astilectron` with `.SetProvisioner(p Provisioner)`. Or you can use the [bootstrap](https://github.com/asticode/go-astilectron-bootstrap) and the [bundler](https://github.com/asticode/go-astilectron-bundler). Check out the [demo](https://github.com/asticode/go-astilectron-demo) to see how to use them.

Beware when trying to add your own app icon as you'll need 2 icons : one compatible with MacOSX (.icns) and one compatible with the rest (.png for instance).

If no BaseDirectoryPath is provided, it defaults to the executable's directory path.

The majority of methods are synchrone which means that when executing them `go-astilectron` will block until it receives a specific Electron event or until the overall context is cancelled. This is the case of `.Start()` which will block until it receives the `app.event.ready` `astilectron` event or until the overall context is cancelled.

## Create a window

```go
// Create a new window
var w, _ = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
    Center: astilectron.PtrBool(true),
    Height: astilectron.PtrInt(600),
    Width:  astilectron.PtrInt(600),
})
w.Create()
```
    
When creating a window you need to indicate a URL as well as options such as position, size, etc.

This is pretty straightforward except the `astilectron.Ptr*` methods so let me explain: GO doesn't do optional fields when json encoding unless you use pointers whereas Electron does handle optional fields. Therefore I added helper methods to convert int, bool and string into pointers and used pointers in structs sent to Electron.

## Open the dev tools

When developing in JS, it's very convenient to debug your code using the browser window's dev tools:

````go
// Open dev tools
w.OpenDevTools()

// Close dev tools
w.CloseDevTools()
````

## Add listeners

```go
// Add a listener on Astilectron
a.On(astilectron.EventNameAppCrash, func(e astilectron.Event) (deleteListener bool) {
    astilog.Error("App has crashed")
    return
})

// Add a listener on the window
w.On(astilectron.EventNameWindowEventResize, func(e astilectron.Event) (deleteListener bool) {
    astilog.Info("Window resized")
    return
})
```
    
Nothing much to say here either except that you can add listeners to Astilectron as well.

## Play with the window

```go
// Play with the window
w.Resize(200, 200)
time.Sleep(time.Second)
w.Maximize()
```
    
Check out the [Window doc](https://godoc.org/github.com/asticode/go-astilectron#Window) for a list of all exported methods

## Send messages from GO to Javascript

### Javascript

```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will listen to messages sent by GO
    astilectron.onMessage(function(message) {
        // Process message
        if (message === "hello") {
            return "world";
        }
    });
})
```

### GO

```go
// This will send a message and execute a callback
// Callbacks are optional
w.SendMessage("hello", func(m *astilectron.EventMessage) {
        // Unmarshal
        var s string
        m.Unmarshal(&s)

        // Process message
        astilog.Debugf("received %s", s)
})
```

This will print `received world` in the GO output

## Send messages from Javascript to GO

### GO

```go
// This will listen to messages sent by Javascript
w.OnMessage(func(m *astilectron.EventMessage) interface{} {
        // Unmarshal
        var s string
        m.Unmarshal(&s)

        // Process message
        if s == "hello" {
                return "world"
        }
        return nil
})
```

### Javascript

```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will send a message to GO
    astilectron.sendMessage("hello", function(message) {
        console.log("received " + message)
    });
})
```

This will print "received world" in the Javascript output

## Play with the window's session

```go
// Clear window's HTTP cache
w.Session.ClearCache()
```

## Handle several screens/displays

```go
// If several displays, move the window to the second display
var displays = a.Displays()
if len(displays) > 1 {
    time.Sleep(time.Second)
    w.MoveInDisplay(displays[1], 50, 50)
}
```

## Menus

```go
// Init a new app menu
// You can do the same thing with a window
var m = a.NewMenu([]*astilectron.MenuItemOptions{
    {
        Label: astilectron.PtrStr("Separator"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Label: astilectron.PtrStr("Normal 1")},
            {
                Label: astilectron.PtrStr("Normal 2"),
                OnClick: func(e astilectron.Event) (deleteListener bool) {
                    astilog.Info("Normal 2 item has been clicked")
                    return
                },
            },
            {Type: astilectron.MenuItemTypeSeparator},
            {Label: astilectron.PtrStr("Normal 3")},
        },
    },
    {
        Label: astilectron.PtrStr("Checkbox"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("Checkbox 1"), Type: astilectron.MenuItemTypeCheckbox},
            {Label: astilectron.PtrStr("Checkbox 2"), Type: astilectron.MenuItemTypeCheckbox},
            {Label: astilectron.PtrStr("Checkbox 3"), Type: astilectron.MenuItemTypeCheckbox},
        },
    },
    {
        Label: astilectron.PtrStr("Radio"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Checked: astilectron.PtrBool(true), Label: astilectron.PtrStr("Radio 1"), Type: astilectron.MenuItemTypeRadio},
            {Label: astilectron.PtrStr("Radio 2"), Type: astilectron.MenuItemTypeRadio},
            {Label: astilectron.PtrStr("Radio 3"), Type: astilectron.MenuItemTypeRadio},
        },
    },
    {
        Label: astilectron.PtrStr("Roles"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Label: astilectron.PtrStr("Minimize"), Role: astilectron.MenuItemRoleMinimize},
            {Label: astilectron.PtrStr("Close"), Role: astilectron.MenuItemRoleClose},
        },
    },
})

// Retrieve a menu item
// This will retrieve the "Checkbox 1" item
mi, _ := m.Item(1, 0)

// Add listener manually
// An OnClick listener has already been added in the options directly for another menu item
mi.On(astilectron.EventNameMenuItemEventClicked, func(e astilectron.Event) bool {
    astilog.Infof("Menu item has been clicked. 'Checked' status is now %t", *e.MenuItemOptions.Checked)
    return false
})

// Create the menu
m.Create()

// Manipulate a menu item
mi.SetChecked(true)

// Init a new menu item
var ni = m.NewItem(&astilectron.MenuItemOptions{
    Label: astilectron.PtrStr("Inserted"),
    SubMenu: []*astilectron.MenuItemOptions{
        {Label: astilectron.PtrStr("Inserted 1")},
        {Label: astilectron.PtrStr("Inserted 2")},
    },
})

// Insert the menu item at position "1"
m.Insert(1, ni)

// Fetch a sub menu
s, _ := m.SubMenu(0)

// Init a new menu item
ni = s.NewItem(&astilectron.MenuItemOptions{
    Label: astilectron.PtrStr("Appended"),
    SubMenu: []*astilectron.MenuItemOptions{
        {Label: astilectron.PtrStr("Appended 1")},
        {Label: astilectron.PtrStr("Appended 2")},
    },
})

// Append menu item dynamically
s.Append(ni)

// Pop up sub menu as a context menu
s.Popup(&astilectron.MenuPopupOptions{PositionOptions: astilectron.PositionOptions{X: astilectron.PtrInt(50), Y: astilectron.PtrInt(50)}})

// Close popup
s.ClosePopup()

// Destroy the menu
m.Destroy()
```

A few things to know:

* when assigning a role to a menu item, `go-astilectron` won't be able to capture its click event
* on MacOS there's no such thing as a window menu, only app menus therefore my advice is to stick to one global app menu instead of creating separate window menus

## Tray

```go
// New tray
var t = a.NewTray(&astilectron.TrayOptions{
    Image:   astilectron.PtrStr("/path/to/image.png"),
    Tooltip: astilectron.PtrStr("Tray's tooltip"),
})

// New tray menu
var m = t.NewMenu([]*astilectron.MenuItemOptions{
    {
        Label: astilectron.PtrStr("Root 1"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Label: astilectron.PtrStr("Item 1")},
            {Label: astilectron.PtrStr("Item 2")},
            {Type: astilectron.MenuItemTypeSeparator},
            {Label: astilectron.PtrStr("Item 3")},
        },
    },
    {
        Label: astilectron.PtrStr("Root 2"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Label: astilectron.PtrStr("Item 1")},
            {Label: astilectron.PtrStr("Item 2")},
        },
    },
})

// Create the menu
m.Create()

// Create tray
t.Create()

// Change tray's image
time.Sleep(time.Second)
t.SetImage(astilectron.PtrStr("/path/to/image-2.png"))
```

## Notifications

```go
// Create the notification
var n = a.NewNotification(&astilectron.NotificationOptions{
	Body: "My Body",
	HasReply: astilectron.PtrBool(true), // Only MacOSX
	Icon: "/path/to/icon",
	ReplyPlaceholder: "type your reply here", // Only MacOSX
	Title: "My title",
})

// Add listeners
n.On(astilectron.EventNameNotificationEventClicked, func(e astilectron.Event) (deleteListener bool) {
	astilog.Debug("the notification has been clicked!")
	return
})
// Only for MacOSX
n.On(astilectron.EventNameNotificationEventReplied, func(e astilectron.Event) (deleteListener bool) {
	astilog.Debugf("the user has replied to the notification: %s", e.Reply)
	return
})

// Create notification
n.Create()

// Show notification
n.Show()
```

## Dock (MacOSX only)

```go
// Get the dock
var d = a.Dock()

// Hide and show the dock
d.Hide()
d.Show()

// Make the Dock bounce
id, _ := d.Bounce(astilectron.DockBounceTypeCritical)

// Cancel the bounce
d.CancelBounce(id)

// Update badge and icon
d.SetBadge("test")
d.SetIcon("/path/to/icon")

// New dock menu
var m = d.NewMenu([]*astilectron.MenuItemOptions{
    {
        Label: astilectron.PtrStr("Root 1"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Label: astilectron.PtrStr("Item 1")},
            {Label: astilectron.PtrStr("Item 2")},
            {Type: astilectron.MenuItemTypeSeparator},
            {Label: astilectron.PtrStr("Item 3")},
        },
    },
        {
        Label: astilectron.PtrStr("Root 2"),
        SubMenu: []*astilectron.MenuItemOptions{
            {Label: astilectron.PtrStr("Item 1")},
            {Label: astilectron.PtrStr("Item 2")},
        },
    },
})

// Create the menu
m.Create()
```

## Dialogs

### Error box

```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will open the dialog
    astilectron.showErrorBox("My Title", "My content")
})
```

### Message box

```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will open the dialog
    astilectron.showMessageBox({message: "My message", title: "My Title"})
})
```

### Open dialog

```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will open the dialog
    astilectron.showOpenDialog({properties: ['openFile', 'multiSelections'], title: "My Title"}, function(paths) {
        console.log("chosen paths are ", paths)
    })
})
```

### Save dialog

```javascript
// This will wait for the astilectron namespace to be ready
document.addEventListener('astilectron-ready', function() {
    // This will open the dialog
    astilectron.showSaveDialog({title: "My title"}, function(filename) {
        console.log("chosen filename is ", filename)
    })
})
```

## Basic auth

```go
// Listen to login events
w.OnLogin(func(i astilectron.Event) (username, password string, err error) {
	// Process the request and auth info
	if i.Request.Method == "GET" && i.AuthInfo.Scheme == "http://" {
		username = "username"
		password = "password"
	}
    return
})
```

# Features and roadmap

- [x] custom branding (custom app name, app icon, etc.)
- [x] window basic methods (create, show, close, resize, minimize, maximize, ...)
- [x] window basic events (close, blur, focus, unresponsive, crashed, ...)
- [x] remote messaging (messages between GO and Javascript)
- [x] single binary distribution
- [x] multi screens/displays
- [x] menu methods and events (create, insert, append, popup, clicked, ...)
- [x] bootstrap
- [x] dialogs (open or save file, alerts, ...)
- [x] tray
- [x] bundler
- [x] session
- [x] accelerators (shortcuts)
- [x] dock
- [x] notifications
- [ ] loader
- [ ] file methods (drag & drop, ...)
- [ ] clipboard methods
- [ ] power monitor events (suspend, resume, ...)
- [ ] desktop capturer (audio and video)
- [ ] window advanced options (add missing ones)
- [ ] window advanced methods (add missing ones)
- [ ] window advanced events (add missing ones)
- [ ] child windows

# Cheers to

[go-thrust](https://github.com/miketheprogrammer/go-thrust) which is awesome but unfortunately not maintained anymore. It inspired this project.
