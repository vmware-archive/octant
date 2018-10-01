const Summary = {
  type: 'summary',
  title: 'Details',
  sections: [
    {
      type: '_primary',
      data: [
        {
          key: 'Name',
          value: 'nginx',
          type: 'string'
        },
        {
          key: 'Namespace',
          value: 'overview',
          type: 'string'
        }
      ]
    },
    {
      type: 'Network',
      data: [
        {
          key: 'Node',
          value: 'docker-for-desktop',
          type: 'link',
          link: '/api/node/blah'
        },
        {
          key: 'IP',
          value: '10.1.68.108',
          type: 'string'
        },
        {
          key: 'health',
          type: 'json',
          data: {
            status: 'OK',
            lastChecked: 'Fri Sep 28 2018 10:37:34 GMT-0700',
            details: {
              cluster: 'Is not broken',
              demo: 'Is going great'
            }
          }
        }
      ]
    }
  ]
}

export default Summary
