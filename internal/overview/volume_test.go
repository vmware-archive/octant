package overview

import (
	"testing"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/apis/core"
)

func Test_summarizePersistentVolumeClaimVolumeSource(t *testing.T) {
	claim := &core.PersistentVolumeClaimVolumeSource{
		ClaimName: "my-claim",
	}

	section := &content.Section{}

	summarizePersistentVolumeClaimVolumeSource(section, claim)

	expected := &content.Section{}
	expected.AddText("Type", "PersistentVolumeClaim")
	expected.AddLink("Claim Name", "my-claim", "/content/overview/config-and-storage/persistent-volume-claims/my-claim")
	expected.AddText("ReadOnly", "false")

	assert.Equal(t, expected, section)
}
