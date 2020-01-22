package component

// Action is an action that can be performed on a component.
type Action struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Form  Form   `json:"form"`
	Modal bool   `json:"modal"`
}
