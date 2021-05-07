package component

type AlertStatus string

const (
	// AlertStatusError is an error alert.
	AlertStatusError AlertStatus = "error"
	// AlertStatusWarning is a warning alert.
	AlertStatusWarning AlertStatus = "warning"
	// AlertStatusInfo is an info alert.
	AlertStatusInfo AlertStatus = "info"
	// AlertStatusSuccess is a success alert.
	AlertStatusSuccess AlertStatus = "success"
	// AlertStatusNeutral is a neutral alert.
	AlertStatusNeutral AlertStatus = "neutral"
)

type AlertType string

const (
	AlertTypeBanner  AlertType = "banner"
	AlertTypeDefault AlertType = "default"
	AlertTypeLight   AlertType = "light"
)

// Alert is an alert. It can be used in components which support alerts.
type Alert struct {
	Status      AlertStatus  `json:"status"`
	Type        AlertType    `json:"type"`
	Message     string       `json:"message"`
	Closable    bool         `json:"closable"`
	ButtonGroup *ButtonGroup `json:"buttonGroup"`
}

// NewAlert creates an instance of Alert.
func NewAlert(alertStatus AlertStatus, alertType AlertType, message string, closable bool, buttonGroup *ButtonGroup) Alert {
	return Alert{
		Status:      alertStatus,
		Type:        alertType,
		Message:     message,
		ButtonGroup: buttonGroup,
		Closable:    closable,
	}
}
