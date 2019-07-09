package component

import (
	"testing"
)

func TestCard_SetAlert(t *testing.T) {
	card := NewCard("card")

	alert := Alert{
		Type:    AlertTypeError,
		Message: "error",
	}

	card.SetAlert(alert)

	expected := NewCard("card")
	expected.Config.Alert = &alert

	AssertEqual(t, expected, card)
}

func TestCard_AddAction(t *testing.T) {
	card := NewCard("card")

	action := Action{
		Name:  "action",
		Title: "action title",
		Form:  Form{},
	}

	card.AddAction(action)

	expected := NewCard("card")
	expected.Config.Actions = []Action{action}

	AssertEqual(t, expected, card)
}

func TestCard_SetBody(t *testing.T) {
	card := NewCard("card")

	body := NewText("body")

	card.SetBody(body)

	expected := NewCard("card")
	expected.Config.Body = body

	AssertEqual(t, expected, card)
}

func TestCardList_AddCard(t *testing.T) {
	card := NewCard("card")

	cardList := NewCardList("list")
	cardList.AddCard(*card)

	expected := NewCardList("list")
	expected.Config.Cards = []Card{*card}

	AssertEqual(t, expected, cardList)
}
