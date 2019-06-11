package event

type notFound interface {
	NotFound() bool
	Path() string
}

