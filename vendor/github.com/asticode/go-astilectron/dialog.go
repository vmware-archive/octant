package astilectron

// Message box types
const (
	MessageBoxTypeError    = "error"
	MessageBoxTypeInfo     = "info"
	MessageBoxTypeNone     = "none"
	MessageBoxTypeQuestion = "question"
	MessageBoxTypeWarning  = "warning"
)

// MessageBoxOptions represents message box options
// We must use pointers since GO doesn't handle optional fields whereas NodeJS does. Use PtrBool, PtrInt or PtrStr
// to fill the struct
// https://github.com/electron/electron/blob/v1.8.1/docs/api/dialog.md#dialogshowmessageboxbrowserwindow-options-callback
type MessageBoxOptions struct {
	Buttons         []string `json:"buttons,omitempty"`
	CancelID        *int     `json:"cancelId,omitempty"`
	CheckboxChecked *bool    `json:"checkboxChecked,omitempty"`
	CheckboxLabel   string   `json:"checkboxLabel,omitempty"`
	ConfirmID       *int     `json:"confirmId,omitempty"`
	DefaultID       *int     `json:"defaultId,omitempty"`
	Detail          string   `json:"detail,omitempty"`
	Icon            string   `json:"icon,omitempty"`
	Message         string   `json:"message,omitempty"`
	NoLink          *bool    `json:"noLink,omitempty"`
	Title           string   `json:"title,omitempty"`
	Type            string   `json:"type,omitempty"`
}
