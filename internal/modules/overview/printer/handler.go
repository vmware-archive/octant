package printer

// Handler configures handlers for a printer.
type Handler interface {
	Handler(printFunc interface{}) error
}

// AddHandlers adds print handlers to a printer.
func AddHandlers(p Handler) error {
	handlers := []interface{}{
		EventListHandler,
		EventHandler,
		ClusterRoleBindingListHandler,
		ClusterRoleBindingHandler,
		ConfigMapListHandler,
		ConfigMapHandler,
		CronJobListHandler,
		CronJobHandler,
		ClusterRoleListHandler,
		ClusterRoleHandler,
		DaemonSetListHandler,
		DaemonSetHandler,
		DeploymentHandler,
		DeploymentListHandler,
		IngressListHandler,
		IngressHandler,
		JobListHandler,
		JobHandler,
		ReplicaSetHandler,
		ReplicaSetListHandler,
		ReplicationControllerHandler,
		ReplicationControllerListHandler,
		PodHandler,
		PodListHandler,
		PersistentVolumeClaimHandler,
		PersistentVolumeClaimListHandler,
		ServiceAccountListHandler,
		ServiceAccountHandler,
		ServiceHandler,
		ServiceListHandler,
		SecretHandler,
		SecretListHandler,
		StatefulSetHandler,
		StatefulSetListHandler,
		RoleBindingListHandler,
		RoleBindingHandler,
		RoleListHandler,
		RoleHandler,
	}

	for _, handler := range handlers {
		if err := p.Handler(handler); err != nil {
			return err
		}
	}

	return nil
}
