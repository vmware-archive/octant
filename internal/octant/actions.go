package octant

import (
	"fmt"
	"time"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ActionDeleteObject            = "action.octant.dev/deleteObject"
	ActionOverviewCordon          = "action.octant.dev/cordon"
	ActionOverviewUncordon        = "action.octant.dev/uncordon"
	ActionOverviewContainerEditor = "action.octant.dev/containerEditor"
	ActionOverviewCronjob         = "action.octant.dev/cronJob"
	ActionOverviewSuspendCronjob  = "action.octant.dev/suspendCronJob"
	ActionOverviewResumeCronjob   = "action.octant.dev/resumeCronJob"
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

func DeleteObjectConfirmationButton(object runtime.Object) (component.ButtonOption, error) {
	if object == nil {
		return nil, fmt.Errorf("object is nil")
	}
	_, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	confirmationTitle := fmt.Sprintf("Delete %s", kind)
	confirmationBody := fmt.Sprintf("Are you sure you want to delete *%s* **%s**? This action is permanent and cannot be recovered.", kind, accessor.GetName())
	return component.WithButtonConfirmation(confirmationTitle, confirmationBody), nil
}

func DeleteObjectConfirmation(object runtime.Object) (*component.Confirmation, error) {
	if object == nil {
		return nil, fmt.Errorf("object is nil")
	}
	_, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()

	accessor, err := meta.Accessor(object)
	if err != nil {
		return nil, err
	}

	confirmationTitle := fmt.Sprintf("Delete %s", kind)
	confirmationBody := fmt.Sprintf("Are you sure you want to delete *%s* **%s**? This action is permanent and cannot be recovered.", kind, accessor.GetName())

	return &component.Confirmation{
		Title: confirmationTitle,
		Body:  confirmationBody,
	}, nil
}
