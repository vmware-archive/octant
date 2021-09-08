package manifest

import (
	"context"
	"fmt"
	"strings"
	"sync"

	imagev5 "github.com/containers/image/v5/image"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

type ImageManifest struct {
	Manifest      string
	Configuration string
}

type ImageEntry struct {
	ImageName string
	HostOS    string
}

type ManifestConfiguration struct {
	imageCache map[ImageEntry]ImageManifest
	imageLock  sync.Mutex
}

var (
	ManifestManager = *NewManifestConfiguration()
)

func NewManifestConfiguration() *ManifestConfiguration {
	mc := &ManifestConfiguration{
		imageCache: make(map[ImageEntry]ImageManifest),
	}
	return mc
}

func (manifest *ManifestConfiguration) SetManifest(imageEntry ImageEntry, imageManifest ImageManifest) {
	manifest.imageCache[imageEntry] = imageManifest
}

func (manifest *ManifestConfiguration) HasEntry(hostOS, imageName string) bool {
	imageName = parseName(imageName)
	_, ok := manifest.imageCache[ImageEntry{ImageName: imageName, HostOS: hostOS}]
	return ok
}

func parseName(imageName string) string {
	parts := strings.SplitN(imageName, "://", 2) // if format not specified, assume docker
	if len(parts) != 2 {
		imageName = "docker://" + imageName
	}
	return imageName
}

func (manifest *ManifestConfiguration) GetImageManifest(ctx context.Context, hostOS, imageName string) (string, string, error) {
	imageName = parseName(imageName)

	imageEntry := ImageEntry{ImageName: imageName, HostOS: hostOS}
	if _, ok := manifest.imageCache[imageEntry]; ok {
		return manifest.imageCache[imageEntry].Manifest, manifest.imageCache[imageEntry].Configuration, nil
	}

	manifest.imageLock.Lock()
	defer manifest.imageLock.Unlock()

	srcRef, err := alltransports.ParseImageName(imageName)
	if err != nil {
		return "", "", fmt.Errorf("error parsing image name for image %s: %w", imageName, err)
	}

	systemCtx := &types.SystemContext{OSChoice: hostOS}

	imageSrc, err := srcRef.NewImageSource(ctx, systemCtx)
	if err != nil {
		return "", "", fmt.Errorf("error creating image source for image %s: %w", imageName, err)
	}
	defer func() {
		err = imageSrc.Close()
	}()

	rawManifest, _, err := imageSrc.GetManifest(ctx, nil)
	if err != nil {
		return "", "", fmt.Errorf("error getting manifest for image %s: %w", imageName, err)
	}

	image, err := imagev5.FromUnparsedImage(ctx, systemCtx, imagev5.UnparsedInstance(imageSrc, nil))
	if err != nil {
		return "", "", fmt.Errorf("error parsing manifest for image %s: %w", imageName, err)
	}

	rawConfiguration, err := image.OCIConfig(ctx)
	if err != nil {
		return "", "", fmt.Errorf("error getting image config blob for image %s: %w", imageName, err)
	}

	configOutput, err := json.MarshalIndent(rawConfiguration, "", "  ")

	manifest.imageCache[imageEntry] = ImageManifest{string(rawManifest), string(configOutput)}

	return string(rawManifest), string(configOutput), err
}
