package printer

import (
	context "context"
	"runtime"
	"strings"
	"sync"

	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

type ImageManifest struct {
	Manifest      string
	Configuration string
}

type ManifestConfiguration struct {
	imageCache map[string]ImageManifest
	imageLock  sync.Mutex
}

var (
	ManifestManager = *NewManifestConfiguration()
)

func NewManifestConfiguration() *ManifestConfiguration {
	mc := &ManifestConfiguration{}
	return mc
}

func (manifest *ManifestConfiguration) GetImageManifest(ctx context.Context, imageName string) (string, string, error) {
	parts := strings.SplitN(imageName, "://", 2) // if format not specified, assume docker
	if len(parts) != 2 {
		imageName = "docker://" + imageName
	}

	if _, ok := manifest.imageCache[imageName]; ok {
		return manifest.imageCache[imageName].Manifest, manifest.imageCache[imageName].Configuration, nil
	}

	manifest.imageLock.Lock()
	defer manifest.imageLock.Unlock()

	srcRef, err := alltransports.ParseImageName(imageName)
	if err != nil {
		return "", "", errors.Wrapf(err, "error parsing image name %q", imageName)
	}

	systemCtx := &types.SystemContext{}

	if runtime.GOOS == "darwin" {
		systemCtx.OSChoice = "linux" // For MAC OS, only linux is currently supported
	}

	imageSrc, err := srcRef.NewImageSource(ctx, systemCtx)
	if err != nil {
		return "", "", errors.Wrapf(err, "error creating image source %q", imageName)
	}

	rawManifest, _, err := imageSrc.GetManifest(ctx, nil)
	if err != nil {
		return "", "", errors.Wrapf(err, "error getting manifest for %q", imageName)
	}

	image, err := image.FromUnparsedImage(ctx, systemCtx, image.UnparsedInstance(imageSrc, nil))
	if err != nil {
		return "", "", errors.Wrapf(err, "Error parsing manifest for %q", imageName)
	}

	rawConfiguration, err := image.OCIConfig(ctx)
	if err != nil {
		return "", "", errors.Wrapf(err, "Error getting image config blob for %q", imageName)
	}

	configOutput, err := json.MarshalIndent(rawConfiguration, "", "  ")

	if manifest.imageCache == nil {
		manifest.imageCache = make(map[string]ImageManifest)
	}
	manifest.imageCache[imageName] = ImageManifest{string(rawManifest), string(configOutput)}

	return string(rawManifest), string(configOutput), nil
}
