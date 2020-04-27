package octant

import (
	"time"

	"github.com/vmware-tanzu/octant/pkg/action"
)

const (
	ActionDeleteObject            = "action.octant.dev/deleteObject"
	ActionOverviewCordon          = "action.octant.dev/cordon"
	ActionOverviewUncordon        = "action.octant.dev/uncordon"
	ActionOverviewContainerEditor = "action.octant.dev/containerEditor"
	ActionOverviewCronjob         = "action.octant.dev/cronJob"
	ActionOverviewServiceEditor   = "action.octant.dev/serviceEditor"
	ActionDeploymentConfiguration = "action.octant.dev/deploymentConfiguration"
	ActionUpdateObject            = "action.octant.dev/update"
)

func sendAlert(alerter action.Alerter, alertType action.AlertType, message string, expiration *time.Time) {
	alert := action.Alert{
		Type:       alertType,
		Message:    message,
		Expiration: expiration,
	}

	alerter.SendAlert(alert)
}
