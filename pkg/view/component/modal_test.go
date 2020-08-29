package component

import (
	"testing"
)

func TestModal_SetBody(t *testing.T) {
	modal := NewModal(TitleFromString("modal"))
	body := NewText("body")
	modal.SetBody(body)

	expected := NewModal(TitleFromString("modal"))
	expected.Config.Body = body

	AssertEqual(t, expected, modal)
}

func TestModal_SetSize(t *testing.T) {
	tests := []struct{
		name string
		size ModalSize
		expected *Modal
	}{
		{
			name: "small",
			size: ModalSizeSmall,
			expected: &Modal{
				Base:   newBase(TypeModal, TitleFromString("modal")),
				Config: ModalConfig{ModalSize: ModalSizeSmall},
			},
		},
		{
			name: "large",
			size: ModalSizeLarge,
			expected: &Modal{
				Base:   newBase(TypeModal, TitleFromString("modal")),
				Config: ModalConfig{ModalSize: ModalSizeLarge},
			},
		},
		{
			name: "extra large",
			size: ModalSizeExtraLarge,
			expected: &Modal{
				Base:   newBase(TypeModal, TitleFromString("modal")),
				Config: ModalConfig{ModalSize: ModalSizeExtraLarge},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			modal := NewModal(TitleFromString("modal"))
			modal.SetSize(test.size)

			AssertEqual(t, test.expected, modal)
		})
	}
}

func TestModal_Open(t *testing.T) {
	modal := NewModal(TitleFromString("modal"))
	modal.Open()

	expected := NewModal(TitleFromString("modal"))
	expected.Config.Opened = true
	AssertEqual(t, expected, modal)
}
