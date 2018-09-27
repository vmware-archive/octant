const Table = {
  type: 'table',
  title: 'Conditions',
  columns: [
    {
      name: 'Type',
      accessor: 'type',
      type: 'string'
    },
    {
      name: 'Status',
      accessor: 'status',
      type: 'string'
    },
    {
      name: 'Last heart beat',
      accessor: 'last_heartbeat_time',
      type: 'time'
    }
  ],
  rows: [
    {
      type: 'Initialized',
      status: 'True',
      last_heartbeat_time: '',
      last_transition_time: '2 minutes',
      reason: '',
      message: ''
    }
  ]
}

export default Table
