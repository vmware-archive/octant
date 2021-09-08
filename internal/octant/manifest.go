package octant

import (
	"context"
	"fmt"

	"github.com/vmware-tanzu/octant/internal/manifest"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
)

type Manifest struct {
	logger log.Logger
}

var _ action.Dispatcher = (*Manifest)(nil)

func NewManifest(logger log.Logger) *Manifest {
	return &Manifest{logger: logger}
}

func (m *Manifest) ActionName() string {
	return ActionGetManifest
}

func (m *Manifest) Handle(ctx context.Context, _ action.Alerter, payload action.Payload) error {
	m.logger.With("payload", payload).Debugf("received action payload")
	image, err := payload.String("image")
	if err != nil {
		return err
	}
	hostOS, err := payload.String("host")
	if err != nil {
		return err
	}

	_, _, err = manifest.ManifestManager.GetImageManifest(ctx, hostOS, image)
	if err != nil {
		return fmt.Errorf("getting manifest for image %s: %w", image, err)
	}
	return nil
}
