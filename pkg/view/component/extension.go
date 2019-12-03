package component

import (
	"encoding/json"
)

type ExtensionConfig struct {
	Tabs []Component `json:"tabs"`
}

func (e *ExtensionConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Tabs []TypedObject `json:"tabs"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	var tabs []Component

	for _, t := range x.Tabs {
		tab, err := t.ToComponent()
		if err != nil {
			return err
		}
		tabs = append(tabs, tab)
	}

	e.Tabs = tabs

	return nil
}

type Extension struct {
	base

	Config ExtensionConfig `json:"config"`
}

func NewExtension() *Extension {
	return &Extension{
		base: newBase(typeExtension, TitleFromString("Extension")),
	}
}

func (e *Extension) AddTab(tab Component) {
	e.Config.Tabs = append(e.Config.Tabs, tab)
}

type extensionMarshal Extension

func (e *Extension) MarshalJSON() ([]byte, error) {
	m := extensionMarshal(*e)
	m.Metadata.Type = typeExtension
	return json.Marshal(&m)
}

func (e *Extension) GetMetadata() Metadata {
	return e.Metadata
}
