package bootstrap

import (
	"github.com/asticode/go-astilectron"
)

// Options represents options
type Options struct {
	Adapter            AstilectronAdapter
	Asset              Asset
	AssetDir           AssetDir
	AstilectronOptions astilectron.Options
	Debug              bool
	MenuOptions        []*astilectron.MenuItemOptions
	OnWait             OnWait
	ResourcesPath      string
	RestoreAssets      RestoreAssets
	TrayMenuOptions    []*astilectron.MenuItemOptions
	TrayOptions        *astilectron.TrayOptions
	Windows            []*Window
}

// Options to setup and create a new window
type Window struct {
	Adapter        WindowAdapter
	Homepage       string
	MessageHandler MessageHandler
	Options        *astilectron.WindowOptions
}

// AstilectronAdapter is a function that adapts the astilectron instance
type AstilectronAdapter func(w *astilectron.Astilectron)

// Asset is a function that retrieves an asset content namely the go-bindata's Asset method
type Asset func(name string) ([]byte, error)

// AssetDir is a function that retrieves an asset dir namely the go-bindata's AssetDir method
type AssetDir func(name string) ([]string, error)

// MessageHandler is a functions that handles messages
type MessageHandler func(w *astilectron.Window, m MessageIn) (payload interface{}, err error)

// OnWait is a function that executes custom actions before waiting
type OnWait func(a *astilectron.Astilectron, w []*astilectron.Window, m *astilectron.Menu, t *astilectron.Tray, tm *astilectron.Menu) error

// RestoreAssets is a function that restores assets namely the go-bindata's RestoreAssets method
type RestoreAssets func(dir, name string) error

// WindowAdapter is a function that adapts a window
type WindowAdapter func(w *astilectron.Window)
