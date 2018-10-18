package content

// Content is content served by the overview API.
type Content interface {
	IsEmpty() bool
}
