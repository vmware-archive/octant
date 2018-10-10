/*
  ingresses: table (name, age)
  services: table (name, labels, cluster ip, internal endpoints, external endpoints, age)
  pods: table (name, node, status, restarts, age)
 */

export default {
  'api/v1/content/overview/workloads/deployments': {
    contents: [
      {
        type: 'table',
        title: 'Deployments',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Labels', accessor: 'labels' },
          { name: 'Pods', accessor: 'pods' },
          { name: 'Age', accessor: 'age' },
          { name: 'Images', accessor: 'images' }
        ],
        rows: [
          {
            name: 'hq-auth-dbmetrics',
            labels: {
              type: 'array',
              data: ['app:auth', 'component:hq-auth-dbmetrics']
            },
            pods: '1 / 1',
            age: '3 days',
            images: {
              data: ['gcr.io/heptio-prod/hq-auth-dbmetrics:20181004'],
              type: 'array'
            }
          },
          {
            name: 'hq-auth',
            labels: {
              type: 'array',
              data: ['app:auth', 'component:hq-auth']
            },
            pods: '3 / 3',
            age: '3 months',
            images: {
              data: ['gcr.io/heptio-prod/hq-auth:20181004'],
              type: 'array'
            }
          }
        ]
      }
    ]
  },
  'api/v1/content/overview/workloads/pods': {
    contents: [
      {
        type: 'table',
        title: 'Pods',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Node', accessor: 'node' },
          { name: 'Status', accessor: 'status' },
          { name: 'Restarts', accessor: 'restarts' },
          { name: 'Age', accessor: 'age' }
        ],
        rows: [
          {
            name: 'raven-7b56d9ddc5-77kfl',
            node: 'ip-10-50-21-66.us-west-2.compute.internal',
            status: 'Running',
            restarts: 0,
            age: '3 days'
          },
          {
            name: 'raven-7b56d9ddc5-7swq6',
            node: 'ip-10-50-31-196.us-west-2.compute.internal',
            status: 'Running',
            restarts: 0,
            age: '3 days'
          },
          {
            name: 'raven-7b56d9ddc5-ndqlj',
            node: 'ip-10-50-0-48.us-west-2.compute.internal',
            status: 'Running',
            restarts: 0,
            age: '3 days'
          }
        ]
      }
    ]
  },
  'api/v1/content/overview/workloads/replica-sets': {
    contents: [
      {
        type: 'table',
        title: 'Replica Sets',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Labels', accessor: 'labels' },
          { name: 'Pods', accessor: 'pods' },
          { name: 'Age', accessor: 'age' },
          { name: 'Images', accessor: 'images' }
        ],
        rows: [
          {
            name: 'hq-auth-6d977bf5dd',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth',
                'pod-template-hash:2853369188'
              ]
            },
            pods: '3 / 3',
            age: '3 days',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth:20181004']
            }
          },
          {
            name: 'hq-auth-dbmetrics-6b69bbbc59',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth-dbmetrics',
                'pod-template-hash:2625666715'
              ]
            },
            pods: '1 / 1',
            age: '3 days',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth-dbmetrics:20181004']
            }
          },
          {
            name: 'hq-auth-7ccc847749',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth',
                'pod-template-hash:3777403305'
              ]
            },
            pods: '0 / 0',
            age: '20 days',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth:20180918']
            }
          },
          {
            name: 'hq-auth-55cf959957',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth',
                'pod-template-hash:1179515513'
              ]
            },
            pods: '0 / 0',
            age: 'a month',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth:20180905']
            }
          }
        ]
      }
    ]
  },
  'api/v1/content/overview/workloads': {
    contents: [
      {
        type: 'table',
        title: 'Deployments',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Labels', accessor: 'labels' },
          { name: 'Pods', accessor: 'pods' },
          { name: 'Age', accessor: 'age' },
          { name: 'Images', accessor: 'images' }
        ],
        rows: [
          {
            name: 'hq-auth-dbmetrics',
            labels: {
              type: 'array',
              data: ['app:auth', 'component:hq-auth-dbmetrics']
            },
            pods: '1 / 1',
            age: '3 days',
            images: {
              data: ['gcr.io/heptio-prod/hq-auth-dbmetrics:20181004'],
              type: 'array'
            }
          },
          {
            name: 'hq-auth',
            labels: {
              type: 'array',
              data: ['app:auth', 'component:hq-auth']
            },
            pods: '3 / 3',
            age: '3 months',
            images: {
              data: ['gcr.io/heptio-prod/hq-auth:20181004'],
              type: 'array'
            }
          }
        ]
      },
      {
        type: 'table',
        title: 'Pods',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Node', accessor: 'node' },
          { name: 'Status', accessor: 'status' },
          { name: 'Restarts', accessor: 'restarts' },
          { name: 'Age', accessor: 'age' }
        ],
        rows: [
          {
            name: 'raven-7b56d9ddc5-77kfl',
            node: 'ip-10-50-21-66.us-west-2.compute.internal',
            status: 'Running',
            restarts: 0,
            age: '3 days'
          },
          {
            name: 'raven-7b56d9ddc5-7swq6',
            node: 'ip-10-50-31-196.us-west-2.compute.internal',
            status: 'Running',
            restarts: 0,
            age: '3 days'
          },
          {
            name: 'raven-7b56d9ddc5-ndqlj',
            node: 'ip-10-50-0-48.us-west-2.compute.internal',
            status: 'Running',
            restarts: 0,
            age: '3 days'
          }
        ]
      },
      {
        type: 'table',
        title: 'Replica Sets',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Labels', accessor: 'labels' },
          { name: 'Pods', accessor: 'pods' },
          { name: 'Age', accessor: 'age' },
          { name: 'Images', accessor: 'images' }
        ],
        rows: [
          {
            name: 'hq-auth-6d977bf5dd',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth',
                'pod-template-hash:2853369188'
              ]
            },
            pods: '3 / 3',
            age: '3 days',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth:20181004']
            }
          },
          {
            name: 'hq-auth-dbmetrics-6b69bbbc59',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth-dbmetrics',
                'pod-template-hash:2625666715'
              ]
            },
            pods: '1 / 1',
            age: '3 days',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth-dbmetrics:20181004']
            }
          },
          {
            name: 'hq-auth-7ccc847749',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth',
                'pod-template-hash:3777403305'
              ]
            },
            pods: '0 / 0',
            age: '20 days',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth:20180918']
            }
          },
          {
            name: 'hq-auth-55cf959957',
            labels: {
              type: 'array',
              data: [
                'app:auth',
                'component: hq-auth',
                'pod-template-hash:1179515513'
              ]
            },
            pods: '0 / 0',
            age: 'a month',
            images: {
              type: 'array',
              data: ['gcr.io/heptio-prod/hq-auth:20180905']
            }
          }
        ]
      }
    ]
  },
  'api/v1/content/overview/discovery-and-load-balancing/ingresses': {
    contents: [
      {
        type: 'table',
        title: 'Ingresses',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Endpoints', accessor: 'endpoints' },
          { name: 'Age', accessor: 'age' }
        ],
        rows: [
          {
            name: 'raven-rest',
            endpoints: '',
            age: '3 months'
          },
          {
            name: 'raven-installer',
            endpoints: '',
            age: '3 months'
          },
          {
            name: 'raven-grpc',
            endpoints: '',
            age: '3 months'
          }
        ]
      }
    ]
  },
  'api/v1/content/overview/discovery-and-load-balancing/services': {
    contents: [
      {
        type: 'table',
        title: 'Services',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Labels', accessor: 'labels' },
          { name: 'Cluster IP', accessor: 'cluster_ip' },
          { name: 'Internal endpoints', accessor: 'internal_endpoints' },
          { name: 'External endpoints', accessor: 'external_endpoints' },
          { name: 'Age', accessor: 'age' }
        ],
        rows: [
          {
            name: 'hq-auth-dbmetrics',
            labels: {
              type: 'array',
              data: ['heptio.com/metrics:metricsz']
            },
            cluster_ip: 'None',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-dbmetrics.heptio-test-auth:7777 TCP',
                'hq-auth-dbmetrics.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 days'
          },
          {
            name: 'hq-auth-rest',
            labels: {
              type: 'array',
              data: []
            },
            cluster_ip: '10.103.143.20',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-dbmetrics.heptio-test-auth:443 TCP',
                'hq-auth-dbmetrics.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 months'
          },
          {
            name: 'hq-auth-grpc',
            labels: {
              type: 'array',
              data: []
            },
            cluster_ip: '10.104.63.139',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-grpc.heptio-test-auth:443 TCP',
                'hq-auth-grpc.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 months'
          },
          {
            name: 'hq-auth-metrics',
            labels: {
              type: 'array',
              data: ['heptio.com/metric:metricz']
            },
            cluster_ip: 'None',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-metrics.heptio-test-auth:443 TCP',
                'hq-auth-metrics.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 months'
          }
        ]
      }
    ]
  },
  'api/v1/content/overview/discovery-and-load-balancing': {
    contents: [
      {
        type: 'table',
        title: 'Ingresses',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Endpoints', accessor: 'endpoints' },
          { name: 'Age', accessor: 'age' }
        ],
        rows: [
          {
            name: 'raven-rest',
            endpoints: '',
            age: '3 months'
          },
          {
            name: 'raven-installer',
            endpoints: '',
            age: '3 months'
          },
          {
            name: 'raven-grpc',
            endpoints: '',
            age: '3 months'
          }
        ]
      },
      {
        type: 'table',
        title: 'Services',
        columns: [
          { name: 'Name', accessor: 'name' },
          { name: 'Labels', accessor: 'labels' },
          { name: 'Cluster IP', accessor: 'cluster_ip' },
          { name: 'Internal endpoints', accessor: 'internal_endpoints' },
          { name: 'External endpoints', accessor: 'external_endpoints' },
          { name: 'Age', accessor: 'age' }
        ],
        rows: [
          {
            name: 'hq-auth-dbmetrics',
            labels: {
              type: 'array',
              data: ['heptio.com/metrics:metricsz']
            },
            cluster_ip: 'None',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-dbmetrics.heptio-test-auth:7777 TCP',
                'hq-auth-dbmetrics.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 days'
          },
          {
            name: 'hq-auth-rest',
            labels: {
              type: 'array',
              data: []
            },
            cluster_ip: '10.103.143.20',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-dbmetrics.heptio-test-auth:443 TCP',
                'hq-auth-dbmetrics.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 months'
          },
          {
            name: 'hq-auth-grpc',
            labels: {
              type: 'array',
              data: []
            },
            cluster_ip: '10.104.63.139',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-grpc.heptio-test-auth:443 TCP',
                'hq-auth-grpc.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 months'
          },
          {
            name: 'hq-auth-metrics',
            labels: {
              type: 'array',
              data: ['heptio.com/metric:metricz']
            },
            cluster_ip: 'None',
            internal_endpoints: {
              type: 'array',
              data: [
                'hq-auth-metrics.heptio-test-auth:443 TCP',
                'hq-auth-metrics.heptio-test-auth:0 TCP'
              ]
            },
            external_endpoints: '',
            age: '3 months'
          }
        ]
      }
    ]
  }
}
