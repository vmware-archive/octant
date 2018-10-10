export default {
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
}
