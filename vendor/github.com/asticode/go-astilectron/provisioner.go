package astilectron

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/asticode/go-astikit"
)

// Var
var (
	regexpDarwinInfoPList = regexp.MustCompile("<string>Electron")
)

// Provisioner represents an object capable of provisioning Astilectron
type Provisioner interface {
	Provision(ctx context.Context, appName, os, arch, versionAstilectron, versionElectron string, p Paths) error
}

// mover is a function that moves a package
type mover func(ctx context.Context, p Paths) error

// defaultProvisioner represents the default provisioner
type defaultProvisioner struct {
	l                astikit.SeverityLogger
	moverAstilectron mover
	moverElectron    mover
}

func newDefaultProvisioner(l astikit.StdLogger) (dp *defaultProvisioner) {
	d := astikit.NewHTTPDownloader(astikit.HTTPDownloaderOptions{
		Sender: astikit.HTTPSenderOptions{
			Logger: l,
		},
	})
	dp = &defaultProvisioner{l: astikit.AdaptStdLogger(l)}
	dp.moverAstilectron = func(ctx context.Context, p Paths) (err error) {
		if err = Download(ctx, dp.l, d, p.AstilectronDownloadSrc(), p.AstilectronDownloadDst()); err != nil {
			return fmt.Errorf("downloading %s into %s failed: %w", p.AstilectronDownloadSrc(), p.AstilectronDownloadDst(), err)
		}
		return
	}
	dp.moverElectron = func(ctx context.Context, p Paths) (err error) {
		if err = Download(ctx, dp.l, d, p.ElectronDownloadSrc(), p.ElectronDownloadDst()); err != nil {
			return fmt.Errorf("downloading %s into %s failed: %w", p.ElectronDownloadSrc(), p.ElectronDownloadDst(), err)
		}
		return
	}
	return
}

// provisionStatusElectronKey returns the electron's provision status key
func provisionStatusElectronKey(os, arch string) string {
	return fmt.Sprintf("%s-%s", os, arch)
}

// Provision implements the provisioner interface
// TODO Package app using electron instead of downloading Electron + Astilectron separately
func (p *defaultProvisioner) Provision(ctx context.Context, appName, os, arch, versionAstilectron, versionElectron string, paths Paths) (err error) {
	// Retrieve provision status
	var s ProvisionStatus
	if s, err = p.ProvisionStatus(paths); err != nil {
		err = fmt.Errorf("retrieving provisioning status failed: %w", err)
		return
	}
	defer p.updateProvisionStatus(paths, &s)

	// Provision astilectron
	if err = p.provisionAstilectron(ctx, paths, s, versionAstilectron); err != nil {
		err = fmt.Errorf("provisioning astilectron failed: %w", err)
		return
	}
	s.Astilectron = &ProvisionStatusPackage{Version: versionAstilectron}

	// Provision electron
	if err = p.provisionElectron(ctx, paths, s, appName, os, arch, versionElectron); err != nil {
		err = fmt.Errorf("provisioning electron failed: %w", err)
		return
	}
	s.Electron[provisionStatusElectronKey(os, arch)] = &ProvisionStatusPackage{Version: versionElectron}
	return
}

// ProvisionStatus represents the provision status
type ProvisionStatus struct {
	Astilectron *ProvisionStatusPackage            `json:"astilectron,omitempty"`
	Electron    map[string]*ProvisionStatusPackage `json:"electron,omitempty"`
}

// ProvisionStatusPackage represents the provision status of a package
type ProvisionStatusPackage struct {
	Version string `json:"version"`
}

// ProvisionStatus returns the provision status
func (p *defaultProvisioner) ProvisionStatus(paths Paths) (s ProvisionStatus, err error) {
	// Open the file
	var f *os.File
	s.Electron = make(map[string]*ProvisionStatusPackage)
	if f, err = os.Open(paths.ProvisionStatus()); err != nil {
		if !os.IsNotExist(err) {
			err = fmt.Errorf("opening file %s failed: %w", paths.ProvisionStatus(), err)
		} else {
			err = nil
		}
		return
	}
	defer f.Close()

	// Unmarshal
	if errLocal := json.NewDecoder(f).Decode(&s); errLocal != nil {
		// For backward compatibility purposes, if there's an unmarshal error we delete the status file and make the
		// assumption that provisioning has to be done all over again
		p.l.Error(fmt.Errorf("json decoding from %s failed: %w", paths.ProvisionStatus(), errLocal))
		p.l.Debugf("Removing %s", f.Name())
		if errLocal = os.RemoveAll(f.Name()); errLocal != nil {
			p.l.Error(fmt.Errorf("removing %s failed: %w", f.Name(), errLocal))
		}
		return
	}
	return
}

// ProvisionStatus updates the provision status
func (p *defaultProvisioner) updateProvisionStatus(paths Paths, s *ProvisionStatus) (err error) {
	// Create the file
	var f *os.File
	if f, err = os.Create(paths.ProvisionStatus()); err != nil {
		err = fmt.Errorf("creating file %s failed: %w", paths.ProvisionStatus(), err)
		return
	}
	defer f.Close()

	// Marshal
	if err = json.NewEncoder(f).Encode(s); err != nil {
		err = fmt.Errorf("json encoding into %s failed: %w", paths.ProvisionStatus(), err)
		return
	}
	return
}

// provisionAstilectron provisions astilectron
func (p *defaultProvisioner) provisionAstilectron(ctx context.Context, paths Paths, s ProvisionStatus, versionAstilectron string) error {
	return p.provisionPackage(ctx, paths, s.Astilectron, p.moverAstilectron, "Astilectron", versionAstilectron, paths.AstilectronUnzipSrc(), paths.AstilectronDirectory(), nil)
}

// provisionElectron provisions electron
func (p *defaultProvisioner) provisionElectron(ctx context.Context, paths Paths, s ProvisionStatus, appName, os, arch, versionElectron string) error {
	return p.provisionPackage(ctx, paths, s.Electron[provisionStatusElectronKey(os, arch)], p.moverElectron, "Electron", versionElectron, paths.ElectronUnzipSrc(), paths.ElectronDirectory(), func() (err error) {
		switch os {
		case "darwin":
			if err = p.provisionElectronFinishDarwin(appName, paths); err != nil {
				return fmt.Errorf("finishing provisioning electron for darwin systems failed: %w", err)
			}
		default:
			p.l.Debug("System doesn't require finshing provisioning electron, moving on...")
		}
		return
	})
}

// provisionPackage provisions a package
func (p *defaultProvisioner) provisionPackage(ctx context.Context, paths Paths, s *ProvisionStatusPackage, m mover, name, version, pathUnzipSrc, pathDirectory string, finish func() error) (err error) {
	// Package has already been provisioned
	if s != nil && s.Version == version {
		p.l.Debugf("%s has already been provisioned to version %s, moving on...", name, version)
		return
	}
	p.l.Debugf("Provisioning %s...", name)

	// Remove previous install
	p.l.Debugf("Removing directory %s", pathDirectory)
	if err = os.RemoveAll(pathDirectory); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("removing %s failed: %w", pathDirectory, err)
	}

	// Move
	if err = m(ctx, paths); err != nil {
		return fmt.Errorf("moving %s failed: %w", name, err)
	}

	// Create directory
	p.l.Debugf("Creating directory %s", pathDirectory)
	if err = os.MkdirAll(pathDirectory, 0755); err != nil {
		return fmt.Errorf("mkdirall %s failed: %w", pathDirectory, err)
	}

	// Unzip
	if err = Unzip(ctx, p.l, pathUnzipSrc, pathDirectory); err != nil {
		return fmt.Errorf("unzipping %s into %s failed: %w", pathUnzipSrc, pathDirectory, err)
	}

	// Finish
	if finish != nil {
		if err = finish(); err != nil {
			return fmt.Errorf("finishing failed: %w", err)
		}
	}
	return
}

// provisionElectronFinishDarwin finishes provisioning electron for Darwin systems
// https://github.com/electron/electron/blob/v1.8.1/docs/tutorial/application-distribution.md#macos
func (p *defaultProvisioner) provisionElectronFinishDarwin(appName string, paths Paths) (err error) {
	// Log
	p.l.Debug("Finishing provisioning electron for darwin system")

	// Custom app icon
	if paths.AppIconDarwinSrc() != "" {
		if err = p.provisionElectronFinishDarwinCopy(paths); err != nil {
			return fmt.Errorf("copying for darwin system finish failed: %w", err)
		}
	}

	// Custom app name
	if appName != "" {
		// Replace
		if err = p.provisionElectronFinishDarwinReplace(appName, paths); err != nil {
			return fmt.Errorf("replacing for darwin system finish failed: %w", err)
		}

		// Rename
		if err = p.provisionElectronFinishDarwinRename(appName, paths); err != nil {
			return fmt.Errorf("renaming for darwin system finish failed: %w", err)
		}
	}
	return
}

// provisionElectronFinishDarwinCopy copies the proper darwin files
func (p *defaultProvisioner) provisionElectronFinishDarwinCopy(paths Paths) (err error) {
	// Icon
	var src, dst = paths.AppIconDarwinSrc(), filepath.Join(paths.ElectronDirectory(), "Electron.app", "Contents", "Resources", "electron.icns")
	if src != "" {
		p.l.Debugf("Copying %s to %s", src, dst)
		if err = astikit.CopyFile(context.Background(), dst, src, astikit.LocalCopyFileFunc); err != nil {
			return fmt.Errorf("copying %s to %s failed: %w", src, dst, err)
		}
	}
	return
}

// provisionElectronFinishDarwinReplace makes the proper replacements in the proper darwin files
func (p *defaultProvisioner) provisionElectronFinishDarwinReplace(appName string, paths Paths) (err error) {
	for _, path := range []string{
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Info.plist"),
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Frameworks", "Electron Helper.app", "Contents", "Info.plist"),
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Frameworks", "Electron Helper (Renderer).app", "Contents", "Info.plist"),
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Frameworks", "Electron Helper (Plugin).app", "Contents", "Info.plist"),
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Frameworks", "Electron Helper (GPU).app", "Contents", "Info.plist"),
	} {
		// Log
		p.l.Debugf("Replacing in %s", path)

		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		// Read file
		var b []byte
		if b, err = ioutil.ReadFile(path); err != nil {
			return fmt.Errorf("reading %s failed: %w", path, err)
		}

		// Open and truncate file
		var f *os.File
		if f, err = os.Create(path); err != nil {
			return fmt.Errorf("creating %s failed: %w", path, err)
		}
		defer f.Close()

		// Replace
		b = regexpDarwinInfoPList.ReplaceAll(b, []byte("<string>"+appName))

		// Write
		if _, err = f.Write(b); err != nil {
			return fmt.Errorf("writing to %s failed: %w", path, err)
		}
	}
	return
}

// rename represents a rename
type rename struct {
	src, dst string
}

// provisionElectronFinishDarwinRename renames the proper darwin folders
func (p *defaultProvisioner) provisionElectronFinishDarwinRename(appName string, paths Paths) (err error) {
	var appDirectory = filepath.Join(paths.electronDirectory, appName+".app")
	var frameworksDirectory = filepath.Join(appDirectory, "Contents", "Frameworks")
	var helper = filepath.Join(frameworksDirectory, appName+" Helper.app")
	var helperRenderer = filepath.Join(frameworksDirectory, appName+" Helper (Renderer).app")
	var helperPlugin = filepath.Join(frameworksDirectory, appName+" Helper (Plugin).app")
	var helperGPU = filepath.Join(frameworksDirectory, appName+" Helper (GPU).app")
	for _, r := range []rename{
		{src: filepath.Join(paths.electronDirectory, "Electron.app"), dst: appDirectory},
		{src: filepath.Join(appDirectory, "Contents", "MacOS", "Electron"), dst: paths.AppExecutable()},
		{src: filepath.Join(frameworksDirectory, "Electron Helper.app"), dst: filepath.Join(helper)},
		{src: filepath.Join(frameworksDirectory, "Electron Helper (Renderer).app"), dst: filepath.Join(helperRenderer)},
		{src: filepath.Join(frameworksDirectory, "Electron Helper (Plugin).app"), dst: filepath.Join(helperPlugin)},
		{src: filepath.Join(frameworksDirectory, "Electron Helper (GPU).app"), dst: filepath.Join(helperGPU)},
		{src: filepath.Join(helper, "Contents", "MacOS", "Electron Helper"), dst: filepath.Join(helper, "Contents", "MacOS", appName+" Helper")},
		{src: filepath.Join(helperRenderer, "Contents", "MacOS", "Electron Helper (Renderer)"), dst: filepath.Join(helperRenderer, "Contents", "MacOS", appName+" Helper (Renderer)")},
		{src: filepath.Join(helperPlugin, "Contents", "MacOS", "Electron Helper (Plugin)"), dst: filepath.Join(helperPlugin, "Contents", "MacOS", appName+" Helper (Plugin)")},
		{src: filepath.Join(helperGPU, "Contents", "MacOS", "Electron Helper (GPU)"), dst: filepath.Join(helperGPU, "Contents", "MacOS", appName+" Helper (GPU)")},
	} {
		p.l.Debugf("Renaming %s into %s", r.src, r.dst)
		if _, err := os.Stat(r.src); os.IsNotExist(err) {
			continue
		}
		if err = os.Rename(r.src, r.dst); err != nil {
			return fmt.Errorf("renaming %s into %s failed: %w", r.src, r.dst, err)
		}
	}
	return
}

// Disembedder is a functions that allows to disembed data from a path
type Disembedder func(src string) ([]byte, error)

// NewDisembedderProvisioner creates a provisioner that can provision based on embedded data
func NewDisembedderProvisioner(d Disembedder, pathAstilectron, pathElectron string, l astikit.StdLogger) Provisioner {
	dp := &defaultProvisioner{l: astikit.AdaptStdLogger(l)}
	dp.moverAstilectron = func(ctx context.Context, p Paths) (err error) {
		if err = Disembed(ctx, dp.l, d, pathAstilectron, p.AstilectronDownloadDst()); err != nil {
			return fmt.Errorf("disembedding %s into %s failed: %w", pathAstilectron, p.AstilectronDownloadDst(), err)
		}
		return
	}
	dp.moverElectron = func(ctx context.Context, p Paths) (err error) {
		if err = Disembed(ctx, dp.l, d, pathElectron, p.ElectronDownloadDst()); err != nil {
			return fmt.Errorf("disembedding %s into %s failed: %w", pathElectron, p.ElectronDownloadDst(), err)
		}
		return
	}
	return dp
}
