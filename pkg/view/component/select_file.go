package component

import "github.com/vmware-tanzu/octant/internal/util/json"

type SelectFile struct {
	Base
	Config SelectFileConfig `json:"config"`
}

type Layout string

const (
	LayoutHorizontal Layout = "horizontal"
	LayoutVertical   Layout = "vertical"
	LayoutCompact    Layout = "compact"
)

type FileStatus string

const (
	FileStatusSuccess FileStatus = "success"
	FileStatusError   FileStatus = "error"
)

type SelectFileConfig struct {
	Label         string     `json:"label"`
	Multiple      bool       `json:"multiple"`
	Status        FileStatus `json:"status"`
	StatusMessage string     `json:"statusMessage"`
	Layout        Layout     `json:"layout"`
	Action        string     `json:"action,omitempty"`
}

// NewSelectFile creates a new Select File component.
func NewSelectFile(label string, multiple bool, layout Layout, action string) *SelectFile {
	sel := &SelectFile{
		Base: newBase(TypeSelectFile, nil),
		Config: SelectFileConfig{
			Label:    label,
			Multiple: multiple,
			Layout:   layout,
			Action:   action,
		},
	}

	return sel
}

// SetStatus sets the status and status message.
func (sf *SelectFile) SetStatus(status FileStatus, message string) {
	sf.Config.Status = status
	sf.Config.StatusMessage = message
}

// GetMetadata accesses the components metadata. Implements Component.
func (sf *SelectFile) GetMetadata() Metadata {
	return sf.Metadata
}

type selectFileMarshal SelectFile

// MarshalJSON implements json.Marshaler.
func (sf *SelectFile) MarshalJSON() ([]byte, error) {
	m := selectFileMarshal(*sf)
	m.Metadata.Type = TypeSelectFile
	m.Metadata.Title = sf.Metadata.Title
	return json.Marshal(&m)
}
