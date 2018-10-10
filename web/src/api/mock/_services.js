export default {
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
}
