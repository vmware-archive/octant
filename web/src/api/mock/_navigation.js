const NavigationData = {
  sections: [
    {
      title: 'Overview',
      path: '/overview',
      children: [
        {
          title: 'Workloads',
          path: '/workloads',
          children: [
            {
              title: 'Cron Jobs',
              path: '/cron-jobs'
            },
            {
              title: 'Daemon Sets',
              path: '/cron-jobs'
            },
            {
              title: 'Deployments',
              path: '/cron-jobs'
            },
            {
              title: 'Jobs',
              path: '/cron-jobs'
            },
            {
              title: 'Pods',
              path: '/cron-jobs'
            },
            {
              title: 'Replica Sets',
              path: '/cron-jobs'
            },
            {
              title: 'Replication Controllers',
              path: '/cron-jobs'
            },
            {
              title: 'Stateful Sets',
              path: '/cron-jobs'
            }
          ]
        },
        {
          title: 'Discovery & Load Balancing',
          path: '/workloads',
          children: [
            {
              title: 'Ingresses',
              path: '/cron-jobs'
            },
            {
              title: 'Services',
              path: '/cron-jobs'
            }
          ]
        },
        {
          title: 'Config & Storage',
          path: '/workloads',
          children: [
            {
              title: 'Config Maps',
              path: '/cron-jobs'
            },
            {
              title: 'Persistent Volume Claims',
              path: '/cron-jobs'
            }
          ]
        },
        {
          title: 'Custom Resources',
          path: '/workloads',
          children: [
            {
              title: 'Certificate',
              path: '/cron-jobs'
            }
          ]
        },
        {
          title: 'RBAC',
          path: '/workloads',
          children: [
            {
              title: 'Roles',
              path: '/cron-jobs'
            },
            {
              title: 'Role Bindings',
              path: '/cron-jobs'
            }
          ]
        }
      ]
    }
  ]
}

export default NavigationData
