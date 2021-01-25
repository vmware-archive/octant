package cluster

// NamespaceInterface is an interface for querying namespace details.
type NamespaceInterface interface {
	Names() ([]string, error)
	InitialNamespace() string
	ProvidedNamespaces() []string
	HasNamespace(namespace string) bool
}
