package overview

// Interface is the Overview interface.
type Interface interface {
	// Namespaces returns a list of namespace names.
	Namespaces() ([]string, error)

	// Navigation returns navigation items for overview.
	Navigation() error

	// Content returns content for a path.
	Content(path string) error
}
