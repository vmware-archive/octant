package astibundler

import (
	"context"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/akavel/rsrc/rsrc"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-bindata"
	"github.com/sam-kamerer/go-plister"
)

// Configuration represents the bundle configuration
type Configuration struct {
	// The app name as it should be displayed everywhere
	// It's also set as an ldflag and therefore accessible in a global var package_name.AppName
	AppName string `json:"app_name"`

	// The bind configuration
	Bind ConfigurationBind `json:"bind"`

	// Whether the app is a darwin agent app
	DarwinAgentApp bool `json:"darwin_agent_app"`

	// List of environments the bundling should be done upon.
	// An environment is a combination of OS and ARCH
	Environments []ConfigurationEnvironment `json:"environments"`

	// The path of the go binary
	// Defaults to "go"
	GoBinaryPath string `json:"go_binary_path"`

	// Paths to icons
	IconPathDarwin  string `json:"icon_path_darwin"` // .icns
	IconPathLinux   string `json:"icon_path_linux"`
	IconPathWindows string `json:"icon_path_windows"` // .ico

	// Info.plist property list
	InfoPlist map[string]interface{} `json:"info_plist"`

	// The path of the project.
	// Defaults to the current directory
	InputPath string `json:"input_path"`

	// Build flags to pass into go build
	BuildFlags map[string]string `json:"build_flags"`

	// LDFlags to pass through to go build
	LDFlags LDFlags `json:"ldflags"`

	// The path used for the LD Flags
	// Defaults to the `Bind.Package` value
	LDFlagsPackage string `json:"ldflags_package"`

	// The path to application manifest file (WINDOWS ONLY)
	ManifestPath string `json:"manifest_path"`

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

	// Show Windows console
	ShowWindowsConsole bool `json:"show_windows_console"`

	// The path where the vendor directory will be created
	// This path must be relative to the output path
	// Defaults to a temp directory
	VendorDirPath string `json:"vendor_dir_path"`

	// Version of Astilectron install
	VersionAstilectron string `json:"version_astilectron"`

	// Version of Electron install
	VersionElectron string `json:"version_electron"`

	// The path to the working directory.
	// Defaults to a temp directory
	WorkingDirectoryPath string `json:"working_directory_path"`

	//!\\ DEBUG ONLY
	AstilectronPath string `json:"astilectron_path"` // when making changes to astilectron
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
	buildFlags           map[string]string
	cancel               context.CancelFunc
	ctx                  context.Context
	d                    *astikit.HTTPDownloader
	darwinAgentApp       bool
	environments         []ConfigurationEnvironment
	infoPlist            map[string]interface{}
	l                    astikit.SeverityLogger
	ldflags              LDFlags
	ldflagsPackage       string
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
	pathManifest         string
	resourcesAdapters    []ConfigurationResourcesAdapter
	showWindowsConsole   bool
	versionAstilectron   string
	versionElectron      string
}

// absPath computes the absolute path
func absPath(configPath string, defaultPathFn func() (string, error)) (o string, err error) {
	if len(configPath) > 0 {
		if o, err = filepath.Abs(configPath); err != nil {
			err = fmt.Errorf("filepath.Abs of %s failed: %w", configPath, err)
			return
		}
	} else if defaultPathFn != nil {
		if o, err = defaultPathFn(); err != nil {
			err = fmt.Errorf("default path function to compute absPath of %s failed: %w", configPath, err)
			return
		}
	}
	return
}

// New builds a new bundler based on a configuration
func New(c *Configuration, l astikit.StdLogger) (b *Bundler, err error) {
	// Init
	b = &Bundler{
		appName:     c.AppName,
		bindPackage: c.Bind.Package,
		d: astikit.NewHTTPDownloader(astikit.HTTPDownloaderOptions{
			Sender: astikit.HTTPSenderOptions{
				Logger: l,
			},
		}),
		environments:       c.Environments,
		darwinAgentApp:     c.DarwinAgentApp,
		resourcesAdapters:  c.ResourcesAdapters,
		l:                  astikit.AdaptStdLogger(l),
		ldflags:            c.LDFlags,
		ldflagsPackage:     c.LDFlagsPackage,
		infoPlist:          c.InfoPlist,
		showWindowsConsole: c.ShowWindowsConsole,
		versionAstilectron: astilectron.DefaultVersionAstilectron,
		versionElectron:    astilectron.DefaultVersionElectron,
	}

	// Ldflags
	if b.ldflags == nil {
		b.ldflags = make(LDFlags)
	}

	if c.VersionAstilectron != "" {
		b.versionAstilectron = c.VersionAstilectron
	}
	if c.VersionElectron != "" {
		b.versionElectron = c.VersionElectron
	}

	if len(c.BuildFlags) > 0 {
		b.buildFlags = c.BuildFlags
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

	// Windows application manifest path
	if b.pathManifest, err = absPath(c.ManifestPath, nil); err != nil {
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

	// Ldflags package
	if len(b.ldflagsPackage) == 0 {
		b.ldflagsPackage = b.bindPackage
	}

	return
}

// HandleSignals handles signals
func (b *Bundler) HandleSignals() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGABRT, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for s := range ch {
			b.l.Infof("Received signal %s", s)
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
	b.l.Debugf("Removing %s", b.pathCache)
	if err = os.RemoveAll(b.pathCache); err != nil {
		err = fmt.Errorf("removing %s failed: %w", b.pathCache, err)
		return
	}
	return
}

// Bundle bundles an astilectron app based on a configuration
func (b *Bundler) Bundle() (err error) {
	// Create the output folder
	if err = os.MkdirAll(b.pathOutput, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed: %w", b.pathOutput, err)
		return
	}

	// Loop through environments
	for _, e := range b.environments {
		b.l.Debugf("Bundling for environment %s/%s", e.OS, e.Arch)
		if err = b.bundle(e); err != nil {
			err = fmt.Errorf("bundling for environment %s/%s failed: %w", e.OS, e.Arch, err)
			return
		}
	}
	return
}

// bundle bundles an os
func (b *Bundler) bundle(e ConfigurationEnvironment) (err error) {
	// Bind data
	b.l.Debug("Binding data")
	if err = b.BindData(e.OS, e.Arch); err != nil {
		err = fmt.Errorf("binding data failed: %w", err)
		return
	}

	// Add windows .syso
	if e.OS == "windows" {
		if err = b.addWindowsSyso(e.Arch); err != nil {
			err = fmt.Errorf("adding windows .syso failed: %w", err)
			return
		}
	}

	// Reset output dir
	var environmentPath = filepath.Join(b.pathOutput, e.OS+"-"+e.Arch)
	if err = b.resetDir(environmentPath); err != nil {
		err = fmt.Errorf("resetting dir %s failed: %w", environmentPath, err)
		return
	}

	std := LDFlags{
		"X": []string{
			b.ldflagsPackage + `.AppName=` + b.appName,
			b.ldflagsPackage + `.BuiltAt=` + time.Now().String(),
			b.ldflagsPackage + `.VersionAstilectron=` + b.versionAstilectron,
			b.ldflagsPackage + `.VersionElectron=` + b.versionElectron,
		},
	}
	if e.OS == "windows" && !b.showWindowsConsole {
		std["H"] = []string{"windowsgui"}
	}
	std.Merge(b.ldflags)

	// Get gopath
	gp := os.Getenv("GOPATH")
	if len(gp) == 0 {
		gp = build.Default.GOPATH
	}

	args := []string{"build", "-ldflags", std.String()}
	var flag string
	for k, v := range b.buildFlags {
		if hasDash := strings.HasPrefix(k, "-"); hasDash {
			flag = k
		} else {
			flag = "-" + k
		}
		if v != "" {
			args = append(args, flag, v)
		} else {
			args = append(args, flag)
		}
	}

	var binaryPath = filepath.Join(environmentPath, "binary")
	args = append(args, "-o", binaryPath, b.pathBuild)

	// Build cmd
	b.l.Debugf("Building for os %s and arch %s astilectron: %s electron: %s", e.OS, e.Arch, b.versionAstilectron, b.versionElectron)
	var cmd = exec.Command(b.pathGoBinary, args...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env,
		"GOARCH="+e.Arch,
		"GOOS="+e.OS,
		"GOPATH="+gp,
	)

	if e.EnvironmentVariables != nil {
		for k, v := range e.EnvironmentVariables {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}

	// Exec
	var o []byte
	b.l.Debugf("Executing %s", strings.Join(cmd.Args, " "))
	if o, err = cmd.CombinedOutput(); err != nil {
		err = fmt.Errorf("building failed: %s", o)
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
	b.l.Debugf("Removing %s", p)
	if err = os.RemoveAll(p); err != nil {
		err = fmt.Errorf("removing %s failed: %w", p, err)
		return
	}

	// Mkdir
	b.l.Debugf("Creating %s", p)
	if err = os.MkdirAll(p, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed: %w", p, err)
		return
	}
	return
}

// BindData binds the data
func (b *Bundler) BindData(os, arch string) (err error) {
	// Reset bind dir
	if err = b.resetDir(b.pathBindInput); err != nil {
		err = fmt.Errorf("resetting dir %s failed: %w", b.pathBindInput, err)
		return
	}

	// Provision the vendor
	if err = b.provisionVendor(os, arch); err != nil {
		err = fmt.Errorf("provisioning the vendor failed: %w", err)
		return
	}

	// Adapt resources
	if err = b.adaptResources(); err != nil {
		err = fmt.Errorf("adapting resources failed: %w", err)
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
	b.l.Debugf("Generating %s", c.Output)
	err = bindata.Translate(c)
	return
}

// provisionVendor provisions the vendor folder
func (b *Bundler) provisionVendor(oS, arch string) (err error) {
	// Create the vendor folder
	b.l.Debugf("Creating %s", b.pathVendor)
	if err = os.MkdirAll(b.pathVendor, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed: %w", b.pathVendor, err)
		return
	}

	// Create the cache folder
	b.l.Debugf("Creating %s", b.pathCache)
	if err = os.MkdirAll(b.pathCache, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed: %w", b.pathCache, err)
		return
	}

	// Provision astilectron
	if err = b.provisionVendorAstilectron(); err != nil {
		err = fmt.Errorf("provisioning astilectron vendor failed: %w", err)
		return
	}

	// Provision electron
	if err = b.provisionVendorElectron(oS, arch); err != nil {
		err = fmt.Errorf("provisioning electron vendor for OS %s and arch %s failed: %w", oS, arch, err)
		return
	}
	return
}

// provisionVendorZip provisions a vendor zip file
func (b *Bundler) provisionVendorZip(pathDownload, pathCache, pathVendor string) (err error) {
	// Download source
	if _, errStat := os.Stat(pathCache); os.IsNotExist(errStat) {
		if err = astilectron.Download(b.ctx, b.l, b.d, pathDownload, pathCache); err != nil {
			err = fmt.Errorf("downloading %s into %s failed: %w", pathDownload, pathCache, err)
			return
		}
	} else {
		b.l.Debugf("%s already exists, skipping download of %s", pathCache, pathDownload)
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}

	// Copy
	b.l.Debugf("Copying %s to %s", pathCache, pathVendor)
	if err = astikit.CopyFile(b.ctx, pathVendor, pathCache, astikit.LocalCopyFileFunc); err != nil {
		err = fmt.Errorf("copying %s to %s failed: %w", pathCache, pathVendor, err)
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

	var p = filepath.Join(b.pathCache, fmt.Sprintf("astilectron-%s.zip", b.versionAstilectron))
	if len(b.pathAstilectron) > 0 {
		// Zip
		b.l.Debugf("Zipping %s into %s", b.pathAstilectron, p)
		if err = astikit.Zip(b.ctx, p+"/"+fmt.Sprintf("astilectron-%s", b.versionAstilectron), b.pathAstilectron); err != nil {
			err = fmt.Errorf("zipping %s into %s failed: %w", b.pathAstilectron, p, err)
			return
		}

		// Check context error
		if b.ctx.Err() != nil {
			return b.ctx.Err()
		}
	}
	return b.provisionVendorZip(astilectron.AstilectronDownloadSrc(b.versionAstilectron), p, filepath.Join(b.pathVendor, zipNameAstilectron))
}

// provisionVendorElectron provisions the electron vendor zip file
func (b *Bundler) provisionVendorElectron(oS, arch string) error {
	return b.provisionVendorZip(
		astilectron.ElectronDownloadSrc(oS, arch, b.versionElectron),
		filepath.Join(b.pathCache, fmt.Sprintf("electron-%s-%s-%s.zip", oS, arch, b.versionElectron)),
		filepath.Join(b.pathVendor, zipNameElectron))
}

func (b *Bundler) adaptResources() (err error) {
	// Create dir
	var o = filepath.Join(b.pathBindInput, b.pathResources)
	b.l.Debugf("Creating %s", o)
	if err = os.MkdirAll(o, 0755); err != nil {
		err = fmt.Errorf("mkdirall %s failed: %w", o, err)
		return
	}

	// Copy resources
	var i = filepath.Join(b.pathInput, b.pathResources)
	b.l.Debugf("Copying %s to %s", i, o)
	if err = astikit.CopyFile(b.ctx, o, i, astikit.LocalCopyFileFunc); err != nil {
		err = fmt.Errorf("copying %s to %s failed: %w", i, o, err)
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
		b.l.Debugf("Running %s in directory %s", strings.Join(cmd.Args, " "), cmd.Dir)
		var b []byte
		if b, err = cmd.CombinedOutput(); err != nil {
			err = fmt.Errorf("running %s failed with output %s", strings.Join(cmd.Args, " "), b)
			return
		}
	}
	return
}

// addWindowsSyso adds the proper windows .syso if needed
func (b *Bundler) addWindowsSyso(arch string) (err error) {
	if len(b.pathIconWindows) > 0 || len(b.pathManifest) > 0 {
		var p = filepath.Join(b.pathInput, "windows.syso")
		b.l.Debugf("Running rsrc for icon %s into %s", b.pathIconWindows, p)
		if err = rsrc.Embed(p, arch, b.pathManifest, b.pathIconWindows); err != nil {
			err = fmt.Errorf("running rsrc for icon %s into %s failed: %w", b.pathIconWindows, p, err)
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
	b.l.Debugf("Creating %s", macOSPath)
	if err = os.MkdirAll(macOSPath, 0777); err != nil {
		err = fmt.Errorf("mkdirall of %s failed: %w", macOSPath, err)
		return
	}

	var macOSBinaryPath = filepath.Join(macOSPath, b.appName)

	var infoPlist *plister.InfoPlist
	if b.infoPlist != nil {
		infoPlist = plister.MapToInfoPlist(b.infoPlist)
		binaryName, _ := infoPlist.Get("CFBundleExecutable").(string)
		if binaryName != "" {
			macOSBinaryPath = filepath.Join(macOSPath, binaryName)
		}
	}

	// Move binary
	b.l.Debugf("Moving %s to %s", binaryPath, macOSBinaryPath)
	if err = astikit.MoveFile(b.ctx, macOSBinaryPath, binaryPath, astikit.LocalCopyFileFunc); err != nil {
		err = fmt.Errorf("moving %s to %s failed: %w", binaryPath, macOSBinaryPath, err)
		return
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}

	// Make sure the binary is executable
	b.l.Debugf("Chmoding %s", macOSBinaryPath)
	if err = os.Chmod(macOSBinaryPath, 0777); err != nil {
		err = fmt.Errorf("chmoding %s failed: %w", macOSBinaryPath, err)
		return
	}

	// Icon
	if len(b.pathIconDarwin) > 0 {
		// Create Resources folder
		var resourcesPath = filepath.Join(contentsPath, "Resources")
		b.l.Debugf("Creating %s", resourcesPath)
		if err = os.MkdirAll(resourcesPath, 0777); err != nil {
			err = fmt.Errorf("mkdirall of %s failed: %w", resourcesPath, err)
			return
		}

		iconFileName := b.appName + filepath.Ext(b.pathIconDarwin)

		if infoPlist != nil {
			ifn, _ := infoPlist.Get("CFBundleIconFile").(string)
			if ifn != "" {
				iconFileName = ifn
			}
		}

		// Copy icon
		var ip = filepath.Join(resourcesPath, iconFileName)
		b.l.Debugf("Copying %s to %s", b.pathIconDarwin, ip)
		if err = astikit.CopyFile(b.ctx, ip, b.pathIconDarwin, astikit.LocalCopyFileFunc); err != nil {
			err = fmt.Errorf("copying %s to %s failed: %w", b.pathIconDarwin, ip, err)
			return
		}

		// Check context error
		if b.ctx.Err() != nil {
			return b.ctx.Err()
		}
	}

	// Add Info.plist file
	var fp = filepath.Join(contentsPath, "Info.plist")
	b.l.Debugf("Adding Info.plist to %s", fp)

	if infoPlist != nil {
		if err = plister.Generate(fp, infoPlist); err != nil {
			err = fmt.Errorf("generating Info.plist failed: %w", err)
		}
		return
	}

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
		err = fmt.Errorf("adding Info.plist to %s failed: %w", fp, err)
		return
	}
	return
}

// finishLinux finishes bundling for a linux system
// TODO Add .desktop file
func (b *Bundler) finishLinux(environmentPath, binaryPath string) (err error) {
	// Move binary
	var linuxBinaryPath = filepath.Join(environmentPath, b.appName)
	b.l.Debugf("Moving %s to %s", binaryPath, linuxBinaryPath)
	if err = astikit.MoveFile(b.ctx, linuxBinaryPath, binaryPath, astikit.LocalCopyFileFunc); err != nil {
		err = fmt.Errorf("moving %s to %s failed: %w", binaryPath, linuxBinaryPath, err)
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
	b.l.Debugf("Moving %s to %s", binaryPath, windowsBinaryPath)
	if err = astikit.MoveFile(b.ctx, windowsBinaryPath, binaryPath, astikit.LocalCopyFileFunc); err != nil {
		err = fmt.Errorf("moving %s to %s failed: %w", binaryPath, windowsBinaryPath, err)
		return
	}

	// Check context error
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}
	return
}
