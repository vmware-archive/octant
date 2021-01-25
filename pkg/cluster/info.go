package cluster

// InfoInterface provides connection details for a cluster
type InfoInterface interface {
	Context() string
	Cluster() string
	Server() string
	User() string
}
