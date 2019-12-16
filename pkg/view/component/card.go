package component

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// CardConfig is configuration for the card component.
type CardConfig struct {
	// Body is the body of the card.
	Body Component `json:"body"`
	// Actions are actions for the card.
	Actions []Action `json:"actions"`
	// Alert is the alert to show for the card.
	Alert *Alert `json:"alert,omitempty"`
}

// UnmarshalJSON unmarshals a card config from JSON.
func (c *CardConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Body    TypedObject `json:"body"`
		Actions []Action    `json:"actions"`
		Alert   *Alert      `json:"alert,omitempty"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	body, err := x.Body.ToComponent()
	if err != nil {
		return err
	}

	c.Body = body
	c.Actions = x.Actions
	c.Alert = x.Alert

	return nil
}

// Card is a card component.
type Card struct {
	base

	Config CardConfig `json:"config"`
}

// NewCard creates a card component.
func NewCard(title []TitleComponent) *Card {
	return &Card{
		base: newBase(typeCard, title),
	}
}

// AddAction adds an action to a card.
func (c *Card) AddAction(action Action) {
	c.Config.Actions = append(c.Config.Actions, action)
}

// SetBody sets the body for the card.
func (c *Card) SetBody(body Component) {
	c.Config.Body = body
}

// SetAlert sets an alert for a card.
func (c *Card) SetAlert(alert Alert) {
	c.Config.Alert = &alert
}

type cardMarshal Card

// MarshalJSON marshals a card to JSON.
func (c *Card) MarshalJSON() ([]byte, error) {
	m := cardMarshal(*c)
	m.Metadata.Type = typeCard
	return json.Marshal(&m)
}

// CardListConfig is configuration for a card list.
type CardListConfig struct {
	// Cards is a slice of cads.
	Cards []Card `json:"cards"`
}

// UnmarshalJSON unmarshals a card list config from JSON.
func (c *CardListConfig) UnmarshalJSON(data []byte) error {
	x := struct {
		Cards []TypedObject `json:"cards"`
	}{}

	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	for _, typedObject := range x.Cards {
		component, err := typedObject.ToComponent()
		if err != nil {
			return err
		}

		card, ok := component.(*Card)
		if !ok {
			return errors.New("item was not a card")
		}

		c.Cards = append(c.Cards, *card)
	}

	return nil
}

// CardList is a component which comprises of a list of cards.
type CardList struct {
	base
	Config CardListConfig `json:"config"`
}

// NewCardList creates a card list component.
func NewCardList(title string) *CardList {
	return &CardList{
		base: newBase(typeCardList, TitleFromString(title)),
	}
}

var _ Component = (*CardList)(nil)

// AddCard adds a card to the list.
func (c *CardList) AddCard(card Card) {
	c.Config.Cards = append(c.Config.Cards, card)
}

type cardListMarshal CardList

// MarshalJSON marshals a card list to JSON.
func (c *CardList) MarshalJSON() ([]byte, error) {
	m := cardListMarshal(*c)
	m.Metadata.Type = typeCardList
	return json.Marshal(&m)
}
