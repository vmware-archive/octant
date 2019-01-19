package astibundler

import (
	"path/filepath"

	"github.com/asticode/go-astilectron"
)

// Constants
const (
	vendorDirectoryName = "vendor_astilectron_bundler"
	zipNameAstilectron  = "astilectron.zip"
	zipNameElectron     = "electron.zip"
)

// NewProvisioner builds the proper disembedder provisioner
func NewProvisioner(disembedFunc func(string) ([]byte, error)) astilectron.Provisioner {
	return astilectron.NewDisembedderProvisioner(disembedFunc, filepath.Join(vendorDirectoryName, zipNameAstilectron), filepath.Join(vendorDirectoryName, zipNameElectron))
}
