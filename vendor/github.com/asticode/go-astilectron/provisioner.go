package astilectron

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/asticode/go-astilog"
	"github.com/asticode/go-astitools/os"
	"github.com/asticode/go-astitools/regexp"
	"github.com/pkg/errors"
)

// Var
var (
	defaultHTTPClient     = &http.Client{}
	regexpDarwinInfoPList = regexp.MustCompile("<string>Electron")
)

// Provisioner represents an object capable of provisioning Astilectron
type Provisioner interface {
	Provision(ctx context.Context, appName, os, arch string, p Paths) error
}

// mover is a function that moves a package
type mover func(ctx context.Context, p Paths) error

// DefaultProvisioner represents the default provisioner
var DefaultProvisioner = &defaultProvisioner{
	moverAstilectron: func(ctx context.Context, p Paths) (err error) {
		if err = Download(ctx, defaultHTTPClient, p.AstilectronDownloadSrc(), p.AstilectronDownloadDst()); err != nil {
			return errors.Wrapf(err, "downloading %s into %s failed", p.AstilectronDownloadSrc(), p.AstilectronDownloadDst())
		}
		return
	},
	moverElectron: func(ctx context.Context, p Paths) (err error) {
		if err = Download(ctx, defaultHTTPClient, p.ElectronDownloadSrc(), p.ElectronDownloadDst()); err != nil {
			return errors.Wrapf(err, "downloading %s into %s failed", p.ElectronDownloadSrc(), p.ElectronDownloadDst())
		}
		return
	},
}

// defaultProvisioner represents the default provisioner
type defaultProvisioner struct {
	moverAstilectron mover
	moverElectron    mover
}

// provisionStatusElectronKey returns the electron's provision status key
func provisionStatusElectronKey(os, arch string) string {
	return fmt.Sprintf("%s-%s", os, arch)
}

// Provision implements the provisioner interface
// TODO Package app using electron instead of downloading Electron + Astilectron separately
func (p *defaultProvisioner) Provision(ctx context.Context, appName, os, arch string, paths Paths) (err error) {
	// Retrieve provision status
	var s ProvisionStatus
	if s, err = p.ProvisionStatus(paths); err != nil {
		err = errors.Wrap(err, "retrieving provisioning status failed")
		return
	}
	defer p.updateProvisionStatus(paths, &s)

	// Provision astilectron
	if err = p.provisionAstilectron(ctx, paths, s); err != nil {
		err = errors.Wrap(err, "provisioning astilectron failed")
		return
	}
	s.Astilectron = &ProvisionStatusPackage{Version: VersionAstilectron}

	// Provision electron
	if err = p.provisionElectron(ctx, paths, s, appName, os, arch); err != nil {
		err = errors.Wrap(err, "provisioning electron failed")
		return
	}
	s.Electron[provisionStatusElectronKey(os, arch)] = &ProvisionStatusPackage{Version: VersionElectron}
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
			err = errors.Wrapf(err, "opening file %s failed", paths.ProvisionStatus())
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
		astilog.Error(errors.Wrapf(errLocal, "json decoding from %s failed", paths.ProvisionStatus()))
		astilog.Debugf("Removing %s", f.Name())
		if errLocal = os.RemoveAll(f.Name()); errLocal != nil {
			astilog.Error(errors.Wrapf(errLocal, "removing %s failed", f.Name()))
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
		err = errors.Wrapf(err, "creating file %s failed", paths.ProvisionStatus())
		return
	}
	defer f.Close()

	// Marshal
	if err = json.NewEncoder(f).Encode(s); err != nil {
		err = errors.Wrapf(err, "json encoding into %s failed", paths.ProvisionStatus())
		return
	}
	return
}

// provisionAstilectron provisions astilectron
func (p *defaultProvisioner) provisionAstilectron(ctx context.Context, paths Paths, s ProvisionStatus) error {
	return p.provisionPackage(ctx, paths, s.Astilectron, p.moverAstilectron, "Astilectron", VersionAstilectron, paths.AstilectronUnzipSrc(), paths.AstilectronDirectory(), nil)
}

// provisionElectron provisions electron
func (p *defaultProvisioner) provisionElectron(ctx context.Context, paths Paths, s ProvisionStatus, appName, os, arch string) error {
	return p.provisionPackage(ctx, paths, s.Electron[provisionStatusElectronKey(os, arch)], p.moverElectron, "Electron", VersionElectron, paths.ElectronUnzipSrc(), paths.ElectronDirectory(), func() (err error) {
		switch os {
		case "darwin":
			if err = p.provisionElectronFinishDarwin(appName, paths); err != nil {
				return errors.Wrap(err, "finishing provisioning electron for darwin systems failed")
			}
		default:
			astilog.Debug("System doesn't require finshing provisioning electron, moving on...")
		}
		return
	})
}

// provisionPackage provisions a package
func (p *defaultProvisioner) provisionPackage(ctx context.Context, paths Paths, s *ProvisionStatusPackage, m mover, name, version, pathUnzipSrc, pathDirectory string, finish func() error) (err error) {
	// Package has already been provisioned
	if s != nil && s.Version == version {
		astilog.Debugf("%s has already been provisioned to version %s, moving on...", name, version)
		return
	}
	astilog.Debugf("Provisioning %s...", name)

	// Remove previous install
	astilog.Debugf("Removing directory %s", pathDirectory)
	if err = os.RemoveAll(pathDirectory); err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "removing %s failed", pathDirectory)
	}

	// Move
	if err = m(ctx, paths); err != nil {
		return errors.Wrapf(err, "moving %s failed", name)
	}

	// Create directory
	astilog.Debugf("Creating directory %s", pathDirectory)
	if err = os.MkdirAll(pathDirectory, 0755); err != nil {
		return errors.Wrapf(err, "mkdirall %s failed", pathDirectory)
	}

	// Unzip
	if err = Unzip(ctx, pathUnzipSrc, pathDirectory); err != nil {
		return errors.Wrapf(err, "unzipping %s into %s failed", pathUnzipSrc, pathDirectory)
	}

	// Finish
	if finish != nil {
		if err = finish(); err != nil {
			return errors.Wrap(err, "finishing failed")
		}
	}
	return
}

// provisionElectronFinishDarwin finishes provisioning electron for Darwin systems
// https://github.com/electron/electron/blob/v1.8.1/docs/tutorial/application-distribution.md#macos
func (p *defaultProvisioner) provisionElectronFinishDarwin(appName string, paths Paths) (err error) {
	// Log
	astilog.Debug("Finishing provisioning electron for darwin system")

	// Custom app icon
	if paths.AppIconDarwinSrc() != "" {
		if err = p.provisionElectronFinishDarwinCopy(paths); err != nil {
			return errors.Wrap(err, "copying for darwin system finish failed")
		}
	}

	// Custom app name
	if appName != "" {
		// Replace
		if err = p.provisionElectronFinishDarwinReplace(appName, paths); err != nil {
			return errors.Wrap(err, "replacing for darwin system finish failed")
		}

		// Rename
		if err = p.provisionElectronFinishDarwinRename(appName, paths); err != nil {
			return errors.Wrap(err, "renaming for darwin system finish failed")
		}
	}
	return
}

// provisionElectronFinishDarwinCopy copies the proper darwin files
func (p *defaultProvisioner) provisionElectronFinishDarwinCopy(paths Paths) (err error) {
	// Icon
	var src, dst = paths.AppIconDarwinSrc(), filepath.Join(paths.ElectronDirectory(), "Electron.app", "Contents", "Resources", "electron.icns")
	if src != "" {
		astilog.Debugf("Copying %s to %s", src, dst)
		if err = astios.Copy(context.Background(), src, dst); err != nil {
			return errors.Wrapf(err, "copying %s to %s failed", src, dst)
		}
	}
	return
}

// provisionElectronFinishDarwinReplace makes the proper replacements in the proper darwin files
func (p *defaultProvisioner) provisionElectronFinishDarwinReplace(appName string, paths Paths) (err error) {
	for _, p := range []string{
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Info.plist"),
		filepath.Join(paths.electronDirectory, "Electron.app", "Contents", "Frameworks", "Electron Helper.app", "Contents", "Info.plist"),
	} {
		// Log
		astilog.Debugf("Replacing in %s", p)

		// Read file
		var b []byte
		if b, err = ioutil.ReadFile(p); err != nil {
			return errors.Wrapf(err, "reading %s failed", p)
		}

		// Open and truncate file
		var f *os.File
		if f, err = os.Create(p); err != nil {
			return errors.Wrapf(err, "creating %s failed", p)
		}
		defer f.Close()

		// Replace
		astiregexp.ReplaceAll(regexpDarwinInfoPList, &b, []byte("<string>"+appName))

		// Write
		if _, err = f.Write(b); err != nil {
			return errors.Wrapf(err, "writing to %s failed", p)
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
	for _, r := range []rename{
		{src: filepath.Join(paths.electronDirectory, "Electron.app"), dst: appDirectory},
		{src: filepath.Join(appDirectory, "Contents", "MacOS", "Electron"), dst: paths.AppExecutable()},
		{src: filepath.Join(frameworksDirectory, "Electron Helper.app"), dst: filepath.Join(helper)},
		{src: filepath.Join(helper, "Contents", "MacOS", "Electron Helper"), dst: filepath.Join(helper, "Contents", "MacOS", appName+" Helper")},
	} {
		astilog.Debugf("Renaming %s into %s", r.src, r.dst)
		if err = os.Rename(r.src, r.dst); err != nil {
			return errors.Wrapf(err, "renaming %s into %s failed", r.src, r.dst)
		}
	}
	return
}

// Disembedder is a functions that allows to disembed data from a path
type Disembedder func(src string) ([]byte, error)

// NewDisembedderProvisioner creates a provisioner that can provision based on embedded data
func NewDisembedderProvisioner(d Disembedder, pathAstilectron, pathElectron string) Provisioner {
	return &defaultProvisioner{
		moverAstilectron: func(ctx context.Context, p Paths) (err error) {
			if err = Disembed(ctx, d, pathAstilectron, p.AstilectronDownloadDst()); err != nil {
				return errors.Wrapf(err, "disembedding %s into %s failed", pathAstilectron, p.AstilectronDownloadDst())
			}
			return
		},
		moverElectron: func(ctx context.Context, p Paths) (err error) {
			if err = Disembed(ctx, d, pathElectron, p.ElectronDownloadDst()); err != nil {
				return errors.Wrapf(err, "disembedding %s into %s failed", pathElectron, p.ElectronDownloadDst())
			}
			return
		},
	}
}
