export default {
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
}
