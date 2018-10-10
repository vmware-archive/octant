export default {
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
}
