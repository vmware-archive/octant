package overview

var cronJobTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "cron-jobs"),
}

var daemonSetTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "daemon-sets"),
}

var deploymentTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "deployments"),
}

var jobTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "jobs"),
}

var replicaSetTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "replica-sets"),
}

var replicationControllerTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "replication-controllers"),
}

var podTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "pods"),
}

var statefulSetTransforms = map[string]lookupFunc{
	"Name": resourceLink("workloads", "stateful-sets"),
}

var ingressTransforms = map[string]lookupFunc{
	"Name": resourceLink("discovery-and-load-balancing", "ingresses"),
}

var serviceTransforms = map[string]lookupFunc{
	"Name": resourceLink("discovery-and-load-balancing", "services"),
}

var configMapTransforms = map[string]lookupFunc{
	"Name": resourceLink("config-and-storage", "config-maps"),
}

var pvcTransforms = map[string]lookupFunc{
	"Name": resourceLink("config-and-storage", "persistent-volume-claims"),
}

var secretTransforms = map[string]lookupFunc{
	"Name": resourceLink("config-and-storage", "secrets"),
}

var serviceAccountTransforms = map[string]lookupFunc{
	"Name": resourceLink("config-and-storage", "service-accounts"),
}

var roleTransforms = map[string]lookupFunc{
	"Name": resourceLink("rbac", "roles"),
}

var roleBindingTransforms = map[string]lookupFunc{
	"Name": resourceLink("rbac", "role-bindings"),
}
