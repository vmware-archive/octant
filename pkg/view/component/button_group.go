package component

import (
	"fmt"

	"github.com/vmware-tanzu/octant/internal/util/json"
)

// ButtonGroupConfig is configuration for a button group.
type ButtonGroupConfig struct {
	// Buttons are buttons in the group.
	Buttons []Button `json:"buttons"`
}

// ButtonGroup is a group of buttons.
// +octant:component
type ButtonGroup struct {
	Base
	Config ButtonGroupConfig `json:"config"`
}

// NewButtonGroup creates an instance of ButtonGroup.
func NewButtonGroup() *ButtonGroup {
	return &ButtonGroup{
		Base: newBase(TypeButtonGroup, nil),
	}
}

// AddButton adds a button to the ButtonGroup.
func (bg *ButtonGroup) AddButton(button *Button) {
	bg.Config.Buttons = append(bg.Config.Buttons, *button)
}

type buttonGroupMarshal ButtonGroup

// MarshalJSON marshals a button group.
func (bg *ButtonGroup) MarshalJSON() ([]byte, error) {
	m := buttonGroupMarshal(*bg)
	m.Metadata.Type = TypeButtonGroup
	return json.Marshal(&m)
}

func (bg *ButtonGroupConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Buttons []TypedObject `json:"buttons,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, typedObject := range x.Buttons {
		component, err := typedObject.ToComponent()
		if err != nil {
			return err
		}

		button, ok := component.(*Button)
		if !ok {
			return fmt.Errorf("item was not a card")
		}

		bg.Buttons = append(bg.Buttons, *button)
	}

	return nil
}
