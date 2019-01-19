package astibundler

import (
	"context"
	"fmt"
	"go/build"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/akavel/rsrc/rsrc"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/archive"
	"github.com/asticode/go-astitools/os"
	"github.com/asticode/go-bindata"
	"github.com/pkg/errors"
)

// Configuration represents the bundle configuration
type Configuration struct {
	// The app name as it should be displayed everywhere
	// It's also set as an ldflag and therefore accessible in a global var main.AppName
	AppName string `json:"app_name"`

	// The bind configuration
	Bind ConfigurationBind `json:"bind"`

	// Whether the app is a darwin agent app
	DarwinAgentApp bool `json:"darwin_agent_app"`

	// List of environments the bundling should be done upon.
	// An environment is a combination of OS and ARCH
	Environments []ConfigurationEnvironment `json:"environments"`

	// Paths to icons
	IconPathDarwin  string `json:"icon_path_darwin"` // .icns
	IconPathLinux   string `json:"icon_path_linux"`
	IconPathWindows string `json:"icon_path_windows"` // .ico

	// The path of the project.
	// Defaults to the current directory
	InputPath string `json:"input_path"`

	// The path of the go binary
	// Defaults to "go"
	GoBinaryPath string `json:"go_binary_path"`

	// The path where the files will be written
	// Defaults to "output"
	OutputPath string `json:"output_path"`

	// List of commands executed on resources
	// Paths inside commands must be relative to the resources folder
	ResourcesAdapters []ConfigurationResourcesAdapter `json:"resources_adapters"`

	// The path where the resources are/will be created
	// This path must be relative to the input path
	// Defaults to "resources"
	ResourcesPath string `json:"resources_path"`

	// The path where the vendor directory will be created
	// This path must be relative to the output path
	// Defaults to a temp directory
	VendorDirPath string `json:"vendor_dir_path"`

	// The path to the working directory.
	// Defaults to a temp directory
	WorkingDirectoryPath string `json:"working_directory_path"`

	//!\\ DEBUG ONLY
	AstilectronPath string `json:"astilectron_path"` // when making changes to astilectron

	// LDFlags to pass through to go build
	LDFlags LDFlags `json:"ldflags"`
}

type ConfigurationBind struct {
	// The path where the file will be written
	// Defaults to the input path
	OutputPath string `json:"output_path"`

	// The package of the generated file
	// Defaults to "main"
	Package string `json:"package"`
}

// ConfigurationEnvironment represents the bundle configuration environment
type ConfigurationEnvironment struct {
	Arch                 string            `json:"arch"`
	EnvironmentVariables map[string]string `json:"env"`
	OS                   string            `json:"os"`
}

type ConfigurationResourcesAdapter struct {
	Args []string `json:"args"`
	Dir  string   `json:"dir"`
	Name string   `json:"name"`
}

// Bundler represents an object capable of bundling an Astilectron app
type Bundler struct {
	appName              string
	bindPackage          string
	cancel               context.CancelFunc
	Client               *http.Client
	ctx                  context.Context
	darwinAgentApp       bool
	environments         []ConfigurationEnvironment
	ldflags              LDFlags
	pathAstilectron      string
	pathBindInput        string
	pathBindOutput       string
	pathBuild            string
	pathCache            string
	pathIconDarwin       string
	pathIconLinux        string
	pathIconWindows      string
	pathInput            string
	pathGoBinary         string
	pathOutput           string
	pathResources        string
	pathVendor           string
	pathWorkingDirectory string
	resourcesAdapters    []ConfigurationResourcesAdapter
}

// absPath computes the absolute path
func absPath(configPath string, defaultPathFn func() (string, error)) (o string, err error) {
	if len(configPath) > 0 {
		if o, err = filepath.Abs(configPath); err != nil {
			err = errors.Wrapf(err, "filepath.Abs of %s failed", configPath)
			return
		}
	} else if defaultPathFn != nil {
		if o, err = defaultPathFn(); err != nil {
			err = errors.Wrapf(err, "default path function to compute absPath of %s failed", configPath)
			return
		}
	}
	return
}

// New builds a new bundler based on a configuration
func New(c *Configuration) (b *Bundler, err error) {
	// Init
	b = &Bundler{
		appName:           c.AppName,
		bindPackage:       c.Bind.Package,
		Client:            &http.Client{},
		environments:      c.Environments,
		darwinAgentApp:    c.DarwinAgentApp,
		resourcesAdapters: c.ResourcesAdapters,
		ldflags:           c.LDFlags,
	}

	// Add context
	b.ctx, b.cancel = context.WithCancel(context.Background())

	// Loop through environments
	for _, env := range b.environments {
		// Validate OS
		if !astilectron.IsValidOS(env.OS) {
			err = fmt.Errorf("OS %s is invalid", env.OS)
			return
		}
	}

	// Astilectron path
	if b.pathAstilectron, err = absPath(c.AstilectronPath, nil); err != nil {
		return
	}

	// Working directory path
	if b.pathWorkingDirectory, err = absPath(c.WorkingDirectoryPath, func() (string, error) { return filepath.Join(os.TempDir(), "astibundler"), nil }); err != nil {
		return
	}

	// Paths that depend on the working directory path
	b.pathBindInput = filepath.Join(b.pathWorkingDirectory, "bind")
	b.pathCache = filepath.Join(b.pathWorkingDirectory, "cache")

	// Darwin icon path
	if b.pathIconDarwin, err = absPath(c.IconPathDarwin, nil); err != nil {
		return
	}

	// Linux icon path
	if b.pathIconLinux, err = absPath(c.IconPathLinux, nil); err != nil {
		return
	}

	// Windows icon path
	if b.pathIconWindows, err = absPath(c.IconPathWindows, nil); err != nil {
		return
	}

	// Input path
	if b.pathInput, err = absPath(c.InputPath, os.Getwd); err != nil {
		return
	}

	// Paths that depends on the input path
	for _, i := range filepath.SplitList(os.Getenv("GOPATH")) {
		var p = filepath.Join(i, "src")
		if strings.HasPrefix(b.pathInput, p) {
			b.pathBuild = strings.TrimPrefix(strings.TrimPrefix(b.pathInput, p), string(os.PathSeparator))
			break
		}
	}

	// Bind output path
	if b.pathBindOutput, err = absPath(c.Bind.OutputPath, func() (string, error) { return b.pathInput, nil }); err != nil {
		return
	}

	// If build path is empty, ldflags are not set properly
	if len(b.pathBuild) == 0 {
		b.pathBuild = "."
	}

	// Resources path
	if b.pathResources = c.ResourcesPath; len(b.pathResources) == 0 {
		b.pathResources = "resources"
	}

	// Vendor path
	if b.pathVendor = c.VendorDirPath; len(b.pathVendor) == 0 {
		b.pathVendor = vendorDirectoryName
	}
	b.pathVendor = filepath.Join(b.pathBindInput, b.pathVendor)

	// Go binary path
	b.pathGoBinary = "go"
	if len(c.GoBinaryPath) > 0 {
		b.pathGoBinary = c.GoBinaryPath
	}

	// Output path
	if b.pathOutput, err = absPath(c.OutputPath, func() (string, error) {
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}
		p = filepath.Join(p, "output")
		return p, err
	}); err != nil {
		return
	}

	// Bind package
	if len(b.bindPackage) == 0 {
		b.bindPackage = "main"
	}
	return
}

// HandleSignals handles signals
func (b *Bundler) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for s := range ch {
			astilog.Infof("Received signal %s", s)
			b.Stop()
			return
		}
	}()
}

// Stop stops the bundler
func (b *Bundler) Stop() {
	b.cancel()
}

// ClearCache clears the bundler cache
func (b *Bundler) ClearCache() (err error) {
	// Remove cache folder
	astilog.Debugf("Removing %s", b.pathCache)
	if err = os.RemoveAll(b.pathCache); err != nil {
		err = errors.Wrapf(err, "removing %s failed", b.pathCache)
		return
	}
	return
}

// Bundle bundles an astilectron app based on a configuration
func (b *Bundler) Bundle() (err error) {
	// Create the output folder
	astilog.Debugf("Creating %s", b.pathOutput)
	if err = os.MkdirAll(b.pathOutput, 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall %s failed", b.pathOutput)
		return
	}

	// Loop through environments
	for _, e := range b.environments {
		astilog.Debugf("Bundling for environment %s/%s", e.OS, e.Arch)
		if err = b.bundle(e); err != nil {
			err = errors.Wrapf(err, "bundling for environment %s/%s failed", e.OS, e.Arch)
			return
		}
	}
	return
}

// bundle bundles an os
func (b *Bundler) bundle(e ConfigurationEnvironment) (err error) {
	// Bind data
	astilog.Debug("Binding data")
	if err = b.BindData(e.OS, e.Arch); err != nil {
		err = errors.Wrap(err, "binding data failed")
		return
	}

	// Add windows .syso
	if e.OS == "windows" {
		if err = b.addWindowsSyso(e.Arch); err != nil {
			err = errors.Wrap(err, "adding windows .syso failed")
			return
		}
	}

	// Reset output dir
	var environmentPath = filepath.Join(b.pathOutput, e.OS+"-"+e.Arch)
	if err = b.resetDir(environmentPath); err != nil {
		err = errors.Wrapf(err, "resetting dir %s failed", environmentPath)
		return
	}

	std := LDFlags{
		"X": []string{
			`"main.AppName=` + b.appName + `"`,
			`"main.BuiltAt=` + time.Now().String() + `"`,
		},
	}
	if e.OS == "windows" {
		std["H"] = []string{"windowsgui"}
	}

	b.ldflags.merge(std)

	// Get gopath
	gp := os.Getenv("GOPATH")
	if len(gp) == 0 {
		gp = build.Default.GOPATH
	}

	// Build cmd
	astilog.Debugf("Building for os %s and arch %s", e.OS, e.Arch)
	var binaryPath = filepath.Join(environmentPath, "binary")
	var cmd = exec.Command(b.pathGoBinary, "build", "-ldflags", b.ldflags.String(), "-o", binaryPath, b.pathBuild)
	cmd.Env = []string{
		"GOARCH=" + e.Arch,
		"GOOS=" + e.OS,
		"GOCACHE=" + os.Getenv("GOCACHE"),
		"GOFLAGS=" + os.Getenv("GOFLAGS"),
		"GOPATH=" + gp,
		"GOROOT=" + os.Getenv("GOROOT"),
		"PATH=" + os.Getenv("PATH"),
		"TEMP=" + os.Getenv("TEMP"),
		"TAGS=" + os.Getenv("TAGS"),
	}

	if e.EnvironmentVariables != nil {
		for k, v := range e.EnvironmentVariables {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	// Exec
	var o []byte
	astilog.Debugf("Executing %s", strings.Join(cmd.Args, " "))
	if o, err = cmd.CombinedOutput(); err != nil {
		err = errors.Wrapf(err, "building failed: %s", o)
		return
	}

	// Finish bundle based on OS
	switch e.OS {
	case "darwin":
		err = b.finishDarwin(environmentPath, binaryPath)
	case "linux":
		err = b.finishLinux(environmentPath, binaryPath)
	case "windows":
		err = b.finishWindows(environmentPath, binaryPath)
	default:
		err = fmt.Errorf("OS %s is not yet implemented", e.OS)
	}
	return
}

func (b *Bundler) resetDir(p string) (err error) {
	// Remove
	astilog.Debugf("Removing %s", p)
	if err = os.RemoveAll(p); err != nil {
		err = errors.Wrapf(err, "removing %s failed", p)
		return
	}

	// Mkdir
	astilog.Debugf("Creating %s", p)
	if err = os.MkdirAll(p, 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall %s failed", p)
		return
	}
	return
}

// BindData binds the data
func (b *Bundler) BindData(os, arch string) (err error) {
	// Reset bind dir
	if err = b.resetDir(b.pathBindInput); err != nil {
		err = errors.Wrapf(err, "resetting dir %s failed", b.pathBindInput)
		return
	}

	// Provision the vendor
	if err = b.provisionVendor(os, arch); err != nil {
		err = errors.Wrap(err, "provisioning the vendor failed")
		return
	}

	// Adapt resources
	if err = b.adaptResources(); err != nil {
		err = errors.Wrap(err, "adapting resources failed")
		return
	}

	// Build bindata config
	var c = bindata.NewConfig()
	c.Input = []bindata.InputConfig{{Path: b.pathBindInput, Recursive: true}}
	c.Output = filepath.Join(b.pathBindOutput, fmt.Sprintf("bind_%s_%s.go", os, arch))
	c.Package = b.bindPackage
	c.Prefix = b.pathBindInput
	c.Tags = fmt.Sprintf("%s,%s", os, arch)

	// Bind data
	astilog.Debugf("Generating %s", c.Output)
	err = bindata.Translate(c)
	return
}

// provisionVendor provisions the vendor folder
func (b *Bundler) provisionVendor(oS, arch string) (err error) {
	// Create the vendor folder
	astilog.Debugf("Creating %s", b.pathVendor)
	if err = os.MkdirAll(b.pathVendor, 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall %s failed", b.pathVendor)
		return
	}

	// Create the cache folder
	astilog.Debugf("Creating %s", b.pathCache)
	if err = os.MkdirAll(b.pathCache, 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall %s failed", b.pathCache)
		return
	}

	// Provision astilectron
	if err = b.provisionVendorAstilectron(); err != nil {
		err = errors.Wrap(err, "provisioning astilectron vendor failed")
		return
	}

	// Provision electron
	if err = b.provisionVendorElectron(oS, arch); err != nil {
		err = errors.Wrapf(err, "provisioning electron vendor for OS %s and arch %s failed", oS, arch)
		return
	}
	return
}

// provisionVendorZip provisions a vendor zip file
func (b *Bundler) provisionVendorZip(pathDownload, pathCache, pathVendor string) (err error) {
	// Download source
	if _, errStat := os.Stat(pathCache); os.IsNotExist(errStat) {
		if err = astilectron.Download(b.ctx, b.Client, pathDownload, pathCache); err != nil {
			err = errors.Wrapf(err, "downloading %s into %s failed", pathDownload, pathCache)
			return
		}
	} else {
		astilog.Debugf("%s already exists, skipping download of %s", pathCache, pathDownload)
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}

	// Copy
	astilog.Debugf("Copying %s to %s", pathCache, pathVendor)
	if err = astios.Copy(b.ctx, pathCache, pathVendor); err != nil {
		err = errors.Wrapf(err, "copying %s to %s failed", pathCache, pathVendor)
		return
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}
	return
}

// provisionVendorAstilectron provisions the astilectron vendor zip file
func (b *Bundler) provisionVendorAstilectron() (err error) {
	var p = filepath.Join(b.pathCache, fmt.Sprintf("astilectron-%s.zip", astilectron.VersionAstilectron))
	if len(b.pathAstilectron) > 0 {
		// Zip
		astilog.Debugf("Zipping %s into %s", b.pathAstilectron, p)
		if err = astiarchive.Zip(b.ctx, b.pathAstilectron, p, fmt.Sprintf("astilectron-%s", astilectron.VersionAstilectron)); err != nil {
			err = errors.Wrapf(err, "zipping %s into %s failed", b.pathAstilectron, p)
			return
		}

		// Check context error
		if b.ctx.Err() != nil {
			return b.ctx.Err()
		}
	}
	return b.provisionVendorZip(astilectron.AstilectronDownloadSrc(), p, filepath.Join(b.pathVendor, zipNameAstilectron))
}

// provisionVendorElectron provisions the electron vendor zip file
func (b *Bundler) provisionVendorElectron(oS, arch string) error {
	return b.provisionVendorZip(astilectron.ElectronDownloadSrc(oS, arch), filepath.Join(b.pathCache, fmt.Sprintf("electron-%s-%s-%s.zip", oS, arch, astilectron.VersionElectron)), filepath.Join(b.pathVendor, zipNameElectron))
}

func (b *Bundler) adaptResources() (err error) {
	// Create dir
	var o = filepath.Join(b.pathBindInput, b.pathResources)
	astilog.Debugf("Creating %s", o)
	if err = os.MkdirAll(o, 0755); err != nil {
		err = errors.Wrapf(err, "mkdirall %s failed", o)
		return
	}

	// Copy resources
	var i = filepath.Join(b.pathInput, b.pathResources)
	astilog.Debugf("Copying %s to %s", i, o)
	if err = astios.Copy(b.ctx, i, o); err != nil {
		err = errors.Wrapf(err, "copying %s to %s failed", i, o)
		return
	}

	// Nothing to do
	if len(b.resourcesAdapters) == 0 {
		return
	}

	// Loop through adapters
	for _, a := range b.resourcesAdapters {
		// Create cmd
		cmd := exec.CommandContext(b.ctx, a.Name, a.Args...)
		cmd.Dir = o
		if a.Dir != "" {
			cmd.Dir = filepath.Join(o, a.Dir)
		}

		// Run
		astilog.Debugf("Running %s in directory %s", strings.Join(cmd.Args, " "), cmd.Dir)
		var b []byte
		if b, err = cmd.CombinedOutput(); err != nil {
			err = errors.Wrapf(err, "running %s failed with output %s", strings.Join(cmd.Args, " "), b)
			return
		}
	}
	return
}

// addWindowsSyso adds the proper windows .syso if needed
func (b *Bundler) addWindowsSyso(arch string) (err error) {
	if len(b.pathIconWindows) > 0 {
		var p = filepath.Join(b.pathInput, "windows.syso")
		astilog.Debugf("Running rsrc for icon %s into %s", b.pathIconWindows, p)
		if err = rsrc.Embed(p, arch, "", b.pathIconWindows); err != nil {
			err = errors.Wrapf(err, "running rsrc for icon %s into %s failed", b.pathIconWindows, p)
			return
		}
	}
	return
}

// finishDarwin finishes bundling for a darwin system
func (b *Bundler) finishDarwin(environmentPath, binaryPath string) (err error) {
	// Create MacOS folder
	var contentsPath = filepath.Join(environmentPath, b.appName+".app", "Contents")
	var macOSPath = filepath.Join(contentsPath, "MacOS")
	astilog.Debugf("Creating %s", macOSPath)
	if err = os.MkdirAll(macOSPath, 0777); err != nil {
		err = errors.Wrapf(err, "mkdirall of %s failed", macOSPath)
		return
	}

	// Move binary
	var macOSBinaryPath = filepath.Join(macOSPath, b.appName)
	astilog.Debugf("Moving %s to %s", binaryPath, macOSBinaryPath)
	if err = astios.Move(b.ctx, binaryPath, macOSBinaryPath); err != nil {
		err = errors.Wrapf(err, "moving %s to %s failed", binaryPath, macOSBinaryPath)
		return
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}

	// Make sure the binary is executable
	astilog.Debugf("Chmoding %s", macOSBinaryPath)
	if err = os.Chmod(macOSBinaryPath, 0777); err != nil {
		err = errors.Wrapf(err, "chmoding %s failed", macOSBinaryPath)
		return
	}

	// Icon
	if len(b.pathIconDarwin) > 0 {
		// Create Resources folder
		var resourcesPath = filepath.Join(contentsPath, "Resources")
		astilog.Debugf("Creating %s", resourcesPath)
		if err = os.MkdirAll(resourcesPath, 0777); err != nil {
			err = errors.Wrapf(err, "mkdirall of %s failed", resourcesPath)
			return
		}

		// Copy icon
		var ip = filepath.Join(resourcesPath, b.appName+filepath.Ext(b.pathIconDarwin))
		astilog.Debugf("Copying %s to %s", b.pathIconDarwin, ip)
		if err = astios.Copy(b.ctx, b.pathIconDarwin, ip); err != nil {
			err = errors.Wrapf(err, "copying %s to %s failed", b.pathIconDarwin, ip)
			return
		}

		// Check context error
		if b.ctx.Err() != nil {
			return b.ctx.Err()
		}
	}

	// Add Info.plist file
	var fp = filepath.Join(contentsPath, "Info.plist")
	astilog.Debugf("Adding Info.plist to %s", fp)
	lsuiElement := "NO"
	if b.darwinAgentApp {
		lsuiElement = "YES"
	}
	if err = ioutil.WriteFile(fp, []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
	<dict>
		<key>CFBundleIconFile</key>
		<string>`+b.appName+filepath.Ext(b.pathIconDarwin)+`</string>
		<key>CFBundleDisplayName</key>
		<string>`+b.appName+`</string>
		<key>CFBundleExecutable</key>
		<string>`+b.appName+`</string>
		<key>CFBundleName</key>
		<string>`+b.appName+`</string>
		<key>CFBundleIdentifier</key>
		<string>com.`+b.appName+`</string>
		<key>LSUIElement</key>
		<string>`+lsuiElement+`</string>
		<key>CFBundlePackageType</key>
		<string>APPL</string>
	</dict>
</plist>`), 0777); err != nil {
		err = errors.Wrapf(err, "adding Info.plist to %s failed", fp)
		return
	}
	return
}

// finishLinux finishes bundling for a linux system
// TODO Add .desktop file
func (b *Bundler) finishLinux(environmentPath, binaryPath string) (err error) {
	// Move binary
	var linuxBinaryPath = filepath.Join(environmentPath, b.appName)
	astilog.Debugf("Moving %s to %s", binaryPath, linuxBinaryPath)
	if err = astios.Move(b.ctx, binaryPath, linuxBinaryPath); err != nil {
		err = errors.Wrapf(err, "moving %s to %s failed", binaryPath, linuxBinaryPath)
		return
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}
	return
}

// finishWindows finishes bundling for a linux system
func (b *Bundler) finishWindows(environmentPath, binaryPath string) (err error) {
	// Move binary
	var windowsBinaryPath = filepath.Join(environmentPath, b.appName+".exe")
	astilog.Debugf("Moving %s to %s", binaryPath, windowsBinaryPath)
	if err = astios.Move(b.ctx, binaryPath, windowsBinaryPath); err != nil {
		err = errors.Wrapf(err, "moving %s to %s failed", binaryPath, windowsBinaryPath)
		return
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}
	return
}
