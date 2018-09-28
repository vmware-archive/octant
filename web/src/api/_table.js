export default {
  type: 'table',
  title: 'Pods',
  columns: [
    {
      name: 'Name',
      accessor: 'name'
    },
    {
      name: 'Node',
      accessor: 'node'
    },
    {
      name: 'Status',
      accessor: 'status'
    },
    {
      name: 'Restarts',
      accessor: 'restarts'
    },
    {
      name: 'Age',
      accessor: 'age'
    }
  ],
  rows: [
    {
      name: 'raven-56dfc56d88-qn296',
      node: 'node0',
      status: 'running',
      restarts: '0',
      age: 'an hour'
    },
    {
      name: 'raven-56dfc56d88-b4fm4',
      node: 'node1',
      status: 'running',
      restarts: '0',
      age: 'an hour'
    },
    {
      name: 'raven-56dfc56d88-9hdc5',
      node: 'node2',
      status: 'running',
      restarts: '0',
      age: 'an hour'
    }
  ]
}
