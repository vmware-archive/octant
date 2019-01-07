export default {
  title: 'Deployment: nginx-deployment',
  viewComponents: [
    {
      metadata: {
        type: 'grid',
        title: 'Detail',
      },
      config: {
        panels: [
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 0,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'summary',
                  title: 'Details',
                },
                config: {
                  sections: [
                    {
                      header: 'Strategy',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Strategy',
                        },
                        config: {
                          value: 'Rolling Update',
                        },
                      },
                    },
                    {
                      header: 'Selector',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Selector',
                        },
                        config: {
                          value: 'app=nginx-deploy',
                        },
                      },
                    },
                    {
                      header: 'Min Ready Seconds',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Min Ready Seconds',
                        },
                        config: {
                          value: '0',
                        },
                      },
                    },
                    {
                      header: 'Revision History Limit',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Revision History Limit',
                        },
                        config: {
                          value: '10',
                        },
                      },
                    },
                    {
                      header: 'Rolling Update Strategy',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Rolling Update Strategy',
                        },
                        config: {
                          value: 'Max Surge: 25%, Max Unavailable: 25%',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 12,
                y: 0,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'summary',
                  title: 'Details',
                },
                config: {
                  sections: [
                    {
                      content: {
                        header: 'Strategy',
                        metadata: {
                          type: 'text',
                          title: 'Strategy',
                        },
                        config: {
                          value: 'Rolling Update',
                        },
                      },
                    },
                    {
                      header: 'Selector',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Selector',
                        },
                        config: {
                          value: 'app=nginx-deploy',
                        },
                      },
                    },
                    {
                      header: 'Min Ready Seconds',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Min Ready Seconds',
                        },
                        config: {
                          value: '0',
                        },
                      },
                    },
                    {
                      header: 'Revision History Limit',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Revision History Limit',
                        },
                        config: {
                          value: '10',
                        },
                      },
                    },
                    {
                      header: 'Rolling Update Strategy',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Rolling Update Strategy',
                        },
                        config: {
                          value: 'Max Surge: 25%, Max Unavailable: 25%',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 7,
                w: 8,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'quadrant',
                  title: 'Status',
                },
                config: {
                  nw: {
                    value: 1,
                    label: 'Total',
                  },
                  ne: {
                    value: 1,
                    label: 'Updated',
                  },
                  sw: {
                    value: 1,
                    label: 'Available',
                  },
                  se: {
                    value: 1,
                    label: 'Unavailable',
                  },
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 16,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'summary',
                  title: 'Container container-1',
                },
                config: {
                  sections: [
                    {
                      header: 'Image',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Image',
                        },
                        config: {
                          value: 'nginx:1.15',
                        },
                      },
                    },
                    {
                      header: 'Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Port',
                        },
                        config: {
                          value: '80/TCP',
                        },
                      },
                    },
                    {
                      header: 'Host Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Host Port',
                        },
                        config: {
                          value: '0/TCP',
                        },
                      },
                    },
                    {
                      header: 'Environment',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Environment',
                        },
                        config: {
                          value: 'none',
                        },
                      },
                    },
                    {
                      header: 'Mounts',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Mounts',
                        },
                        config: {
                          value: '/usr/share/nginx/html=www(rw)',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 12,
                y: 16,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'summary',
                  title: 'Container container-2',
                },
                config: {
                  sections: [
                    {
                      header: 'Image',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Image',
                        },
                        config: {
                          value: 'nginx:1.15',
                        },
                      },
                    },
                    {
                      header: 'Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Port',
                        },
                        config: {
                          value: '80/TCP',
                        },
                      },
                    },
                    {
                      header: 'Host Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Host Port',
                        },
                        config: {
                          value: '0/TCP',
                        },
                      },
                    },
                    {
                      header: 'Environment',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Environment',
                        },
                        config: {
                          value: 'none',
                        },
                      },
                    },
                    {
                      header: 'Mounts',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Mounts',
                        },
                        config: {
                          value: '/usr/share/nginx/html=www(rw)',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 23,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'summary',
                  title: 'Container container-n',
                },
                config: {
                  sections: [
                    {
                      header: 'Image',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Image',
                        },
                        config: {
                          value: 'nginx:1.15',
                        },
                      },
                    },
                    {
                      header: 'Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Port',
                        },
                        config: {
                          value: '80/TCP',
                        },
                      },
                    },
                    {
                      header: 'Host Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Host Port',
                        },
                        config: {
                          value: '0/TCP',
                        },
                      },
                    },
                    {
                      header: 'Environment',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Environment',
                        },
                        config: {
                          value: 'none',
                        },
                      },
                    },
                    {
                      header: 'Mounts',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Mounts',
                        },
                        config: {
                          value: '/usr/share/nginx/html=www(rw)',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 30,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'table',
                  title: 'Volumes',
                },
                config: {
                  columns: [
                    {
                      name: 'Name',
                      accessor: 'Name',
                    },
                    {
                      name: 'Labels',
                      accessor: 'Labels',
                    },
                    {
                      name: 'Desired',
                      accessor: 'Desired',
                    },
                    {
                      name: 'Current',
                      accessor: 'Current',
                    },
                    {
                      name: 'Ready',
                      accessor: 'Ready',
                    },
                    {
                      name: 'Age',
                      accessor: 'Age',
                    },
                    {
                      name: 'Containers',
                      accessor: 'Containers',
                    },
                    {
                      name: 'Images',
                      accessor: 'Images',
                    },
                    {
                      name: 'Selector',
                      accessor: 'Selector',
                    },
                  ],
                  rows: [
                    {
                      Name: {
                        config: {
                          ref: '/content/overview/workloads/replica-sets/nginx-deployment-7778c58546',
                          value: 'nginx-deployment-7778c58546',
                        },
                        metadata: {
                          type: 'link',
                        },
                      },
                      Age: {
                        config: {
                          timestamp: '1545166466',
                        },
                        metadata: {
                          type: 'timestamp',
                        },
                      },
                      Containers: {
                        config: {
                          value: 'nginx',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Current: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Desired: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Images: {
                        config: {
                          value: 'nginx:1.12',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Labels: {
                        config: {
                          labels: {
                            app: 'nginx-deploy',
                          },
                        },
                        metadata: {
                          type: 'labels',
                        },
                      },
                      Ready: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Selector: {
                        config: {
                          value: 'app=nginx-deploy',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 12,
                y: 30,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'summary',
                  title: 'Additional properties',
                },
                config: {
                  sections: [
                    {
                      header: 'Image',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Image',
                        },
                        config: {
                          value: 'nginx:1.15',
                        },
                      },
                    },
                    {
                      header: 'Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Port',
                        },
                        config: {
                          value: '80/TCP',
                        },
                      },
                    },
                    {
                      header: 'Host Port',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Host Port',
                        },
                        config: {
                          value: '0/TCP',
                        },
                      },
                    },
                    {
                      header: 'Environment',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Environment',
                        },
                        config: {
                          value: 'none',
                        },
                      },
                    },
                    {
                      header: 'Mounts',
                      content: {
                        metadata: {
                          type: 'text',
                          title: 'Mounts',
                        },
                        config: {
                          value: '/usr/share/nginx/html=www(rw)',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 38,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'table',
                  title: 'Conditions',
                },
                config: {
                  columns: [
                    {
                      name: 'Name',
                      accessor: 'Name',
                    },
                    {
                      name: 'Labels',
                      accessor: 'Labels',
                    },
                    {
                      name: 'Desired',
                      accessor: 'Desired',
                    },
                    {
                      name: 'Current',
                      accessor: 'Current',
                    },
                    {
                      name: 'Ready',
                      accessor: 'Ready',
                    },
                    {
                      name: 'Age',
                      accessor: 'Age',
                    },
                    {
                      name: 'Containers',
                      accessor: 'Containers',
                    },
                    {
                      name: 'Images',
                      accessor: 'Images',
                    },
                    {
                      name: 'Selector',
                      accessor: 'Selector',
                    },
                  ],
                  rows: [
                    {
                      Name: {
                        config: {
                          ref: '/content/overview/workloads/replica-sets/nginx-deployment-7778c58546',
                          value: 'nginx-deployment-7778c58546',
                        },
                        metadata: {
                          type: 'link',
                        },
                      },
                      Age: {
                        config: {
                          timestamp: '1545166466',
                        },
                        metadata: {
                          type: 'timestamp',
                        },
                      },
                      Containers: {
                        config: {
                          value: 'nginx',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Current: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Desired: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Images: {
                        config: {
                          value: 'nginx:1.12',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Labels: {
                        config: {
                          labels: {
                            app: 'nginx-deploy',
                          },
                        },
                        metadata: {
                          type: 'labels',
                        },
                      },
                      Ready: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Selector: {
                        config: {
                          value: 'app=nginx-deploy',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 12,
                y: 38,
                w: 12,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'table',
                  title: 'Conditions',
                },
                config: {
                  columns: [
                    {
                      name: 'Name',
                      accessor: 'Name',
                    },
                    {
                      name: 'Labels',
                      accessor: 'Labels',
                    },
                    {
                      name: 'Desired',
                      accessor: 'Desired',
                    },
                    {
                      name: 'Current',
                      accessor: 'Current',
                    },
                    {
                      name: 'Ready',
                      accessor: 'Ready',
                    },
                    {
                      name: 'Age',
                      accessor: 'Age',
                    },
                    {
                      name: 'Containers',
                      accessor: 'Containers',
                    },
                    {
                      name: 'Images',
                      accessor: 'Images',
                    },
                    {
                      name: 'Selector',
                      accessor: 'Selector',
                    },
                  ],
                  rows: [
                    {
                      Name: {
                        config: {
                          ref: '/content/overview/workloads/replica-sets/nginx-deployment-7778c58546',
                          value: 'nginx-deployment-7778c58546',
                        },
                        metadata: {
                          type: 'link',
                        },
                      },
                      Age: {
                        config: {
                          timestamp: '1545166466',
                        },
                        metadata: {
                          type: 'timestamp',
                        },
                      },
                      Containers: {
                        config: {
                          value: 'nginx',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Current: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Desired: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Images: {
                        config: {
                          value: 'nginx:1.12',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Labels: {
                        config: {
                          labels: {
                            app: 'nginx-deploy',
                          },
                        },
                        metadata: {
                          type: 'labels',
                        },
                      },
                      Ready: {
                        config: {
                          value: '3',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                      Selector: {
                        config: {
                          value: 'app=nginx-deploy',
                        },
                        metadata: {
                          type: 'text',
                        },
                      },
                    },
                  ],
                },
              },
            },
          },
          {
            metadata: {
              type: 'panel',
            },
            config: {
              position: {
                x: 0,
                y: 45,
                w: 24,
                h: 7,
              },
              content: {
                metadata: {
                  type: 'table',
                  title: 'Conditions',
                },
                config: {
                  columns: [
                    {
                      name: 'Name',
                      accessor: 'Name',
                    },
                    {
                      name: 'Labels',
                      accessor: 'Labels',
                    },
                    {
                      name: 'Desired',
                      accessor: 'Desired',
                    },
                    {
                      name: 'Current',
                      accessor: 'Current',
                    },
                    {
                      name: 'Ready',
                      accessor: 'Ready',
                    },
                    {
                      name: 'Age',
                      accessor: 'Age',
                    },
                    {
                      name: 'Containers',
                      accessor: 'Containers',
                    },
                    {
                      name: 'Images',
                      accessor: 'Images',
                    },
                    {
                      name: 'Selector',
                      accessor: 'Selector',
                    },
                  ],
                  empty_content: 'Namespace overview does not contain any events for this Deployment',
                },
              },
            },
          },
        ],
      },
    },
    {
      metadata: {
        type: 'resourceViewer',
        title: 'Resource Viewer',
      },
      config: {
        selected: 'ac833d23-c17e-11e8-9212-025000000001',
        adjacencyList: {
          'ac833d23-c17e-11e8-9212-025000000001': [
            {
              node: 'ae54cf9b-0205-11e9-baec-025000000001',
              edge: 'explicit',
            },
          ],
          'ae54cf9b-0205-11e9-baec-025000000001': [
            {
              node: 'pods-ae54cf9b-0205-11e9-baec-025000000001',
              edge: 'explicit',
            },
          ],
        },
        objects: {
          'ac833d23-c17e-11e8-9212-025000000001': {
            name: 'nginx-deployment',
            apiVersion: 'apps/v1',
            kind: 'Deployment',
            status: 'ok',
          },
          'ae54cf9b-0205-11e9-baec-025000000001': {
            name: 'nginx-deployment-7778c58546',
            apiVersion: 'apps/v1',
            kind: 'ReplicaSet',
            status: 'ok',
          },
          'pods-ae54cf9b-0205-11e9-baec-025000000001': {
            name: 'app=nginx-deploy,pod-template-hash=3334714102',
            apiVersion: 'v1',
            kind: 'Pods',
            status: 'ok',
          },
        },
      },
    },
  ],
}
