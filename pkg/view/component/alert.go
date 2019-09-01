package component

type AlertType string

const (
	// AlertTypeError is an error alert.
	AlertTypeError AlertType = "error"
	// AlertTypeWarning is a warning alert.
	AlertTypeWarning AlertType = "warning"
	// AlertTypeInfo is an info alert.
	AlertTypeInfo AlertType = "info"
	// AlertTypeSuccess is a success alert.
	AlertTypeSuccess AlertType = "success"
)

// Alert is an alert. It can be used in components which support alerts.
type Alert struct {
	Type    AlertType `json:"type"`
	Message string    `json:"message"`
}

// NewAlert creates an instance of Alert.
func NewAlert(alertType AlertType, message string) Alert {
	return Alert{
		Type:    alertType,
		Message: message,
	}
}
