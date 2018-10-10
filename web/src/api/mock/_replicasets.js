export default {
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
}
