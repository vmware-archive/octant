package fake

type ClusterInfo struct {
	ContextVal string
	ClusterVal string
	ServerVal  string
	UserVal    string
}

func (ci ClusterInfo) Context() string {
	return ci.ContextVal
}

func (ci ClusterInfo) Cluster() string {
	return ci.ClusterVal
}

func (ci ClusterInfo) Server() string {
	return ci.ServerVal
}

func (ci ClusterInfo) User() string {
	return ci.UserVal
}
