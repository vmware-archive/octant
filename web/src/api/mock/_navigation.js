export default {
  sections: [
    {
      title: 'Overview',
      path: '/content/overview',
      children: [
        {
          title: 'Workloads',
          path: '/content/overview/workloads',
          children: [
            {
              title: 'Cron Jobs',
              path: '/content/overview/workloads/cron-jobs'
            },
            {
              title: 'Daemon Sets',
              path: '/content/overview/workloads/daemon-sets'
            },
            {
              title: 'Deployments',
              path: '/content/overview/workloads/deployments'
            },
            { title: 'Jobs', path: '/content/overview/workloads/jobs' },
            { title: 'Pods', path: '/content/overview/workloads/pods' },
            {
              title: 'Replica Sets',
              path: '/content/overview/workloads/replica-sets'
            },
            {
              title: 'Replication Controllers',
              path: '/content/overview/workloads/replication-controllers'
            },
            {
              title: 'Stateful Sets',
              path: '/content/overview/workloads/stateful-sets'
            }
          ]
        },
        {
          title: 'Discovery and Load Balancing',
          path: '/content/overview/discovery-and-load-balancing',
          children: [
            {
              title: 'Ingresses',
              path: '/content/overview/discovery-and-load-balancing/ingresses'
            },
            {
              title: 'Services',
              path: '/content/overview/discovery-and-load-balancing/services'
            }
          ]
        },
        {
          title: 'Config and Storage',
          path: '/content/overview/config-and-storage',
          children: [
            {
              title: 'Config Maps',
              path: '/content/overview/config-and-storage/config-maps'
            },
            {
              title: 'Persistent Volume Claims',
              path:
                '/content/overview/config-and-storage/persistent-volume-claims'
            },
            {
              title: 'Secrets',
              path: '/content/overview/config-and-storage/secrets'
            }
          ]
        },
        {
          title: 'Custom Resources',
          path: '/content/overview/custom-resources'
        },
        {
          title: 'RBAC',
          path: '/content/overview/rbac',
          children: [
            { title: 'Roles', path: '/content/overview/rbac/roles' },
            {
              title: 'Role Bindings',
              path: '/content/overview/rbac/role-bindings'
            }
          ]
        },
        { title: 'Events', path: '/content/overview/events' }
      ]
    }
  ]
}
