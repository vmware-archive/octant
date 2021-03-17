interface NodeDataDef {
  nodes: any;
  edges?: any;
}

export const REAL_DATA_STATEFUL_SET: NodeDataDef = {
  edges: {
    'kafka pods': [
      { node: '0bf159aa-01ea-4742-a6a2-becef1178827', edge: 'explicit' },
      { node: '14eda8ed-87c3-4aa1-a3cb-9f4279704fc5', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '9cd5e4f0-4979-4119-9c93-7df18bd88059', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
  },
  nodes: {
    '0bf159aa-01ea-4742-a6a2-becef1178827': {
      name: 'kafka-config',
      apiVersion: 'v1',
      kind: 'ConfigMap',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 ConfigMap is OK' } },
      ],
      path: {
        config: {
          value: 'kafka-config',
          ref:
            '/overview/namespace/milan/config-and-storage/config-maps/kafka-config',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '14eda8ed-87c3-4aa1-a3cb-9f4279704fc5': {
      name: 'kafka-headless',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'kafka-headless',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/kafka-headless',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '3c81e771-d723-403d-a19b-be7ce87ff7f2': {
      name: 'default-token-4dln7',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'default-token-4dln7',
          ref:
            '/overview/namespace/milan/config-and-storage/secrets/default-token-4dln7',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '9cd5e4f0-4979-4119-9c93-7df18bd88059': {
      name: 'kafka',
      apiVersion: 'apps/v1',
      kind: 'StatefulSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Stateful Set is OK' } },
      ],
      path: {
        config: {
          value: 'kafka',
          ref: '/overview/namespace/milan/workloads/stateful-sets/kafka',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'a4e5517e-0563-4158-88d3-a0492fe18cd5': {
      name: 'default',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'default',
          ref:
            '/overview/namespace/milan/config-and-storage/service-accounts/default',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'kafka pods': {
      name: 'kafka pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'kafka-0': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
  },
};

export const REAL_DATA_DAEMON_SET: NodeDataDef = {
  edges: {
    'hubble pods': [
      { node: '2ca0da85-f263-4087-a732-73e5501c0a47', edge: 'explicit' },
      { node: 'f5beb4cb-2c7b-474d-9719-0ac02fd8b8b7', edge: 'explicit' },
      { node: 'f69e1b15-a257-42fa-a367-d4a5eb70d8cf', edge: 'explicit' },
      { node: 'f93575bb-0f33-4aa0-8d64-6ebb1cbdf7ce', edge: 'explicit' },
    ],
  },
  nodes: {
    '2ca0da85-f263-4087-a732-73e5501c0a47': {
      name: 'hubble',
      apiVersion: 'apps/v1',
      kind: 'DaemonSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Daemon Set is OK' } },
      ],
      path: {
        config: {
          value: 'hubble',
          ref: '/overview/namespace/kube-system/workloads/daemon-sets/hubble',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'f5beb4cb-2c7b-474d-9719-0ac02fd8b8b7': {
      name: 'hubble-token-smc5q',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'hubble-token-smc5q',
          ref:
            '/overview/namespace/kube-system/config-and-storage/secrets/hubble-token-smc5q',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'f69e1b15-a257-42fa-a367-d4a5eb70d8cf': {
      name: 'hubble-grpc',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Service has no endpoint addresses' },
        },
      ],
      path: {
        config: {
          value: 'hubble-grpc',
          ref:
            '/overview/namespace/kube-system/discovery-and-load-balancing/services/hubble-grpc',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'f93575bb-0f33-4aa0-8d64-6ebb1cbdf7ce': {
      name: 'hubble',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'hubble',
          ref:
            '/overview/namespace/kube-system/config-and-storage/service-accounts/hubble',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'hubble pods': {
      name: 'hubble pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'hubble-4gnq8': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-8lfqv': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-dwzx5': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-jjsdm': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-khckr': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-lhz85': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-q5nwq': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-rqmxz': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-rxbkz': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'hubble-vldtj': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
  },
};

export const REAL_DATA_DAEMON_SET2: NodeDataDef = {
  nodes: {
    '16428c94-a848-47d5-b1e3-c8245b57959b': {
      name: 'metadata-proxy-v0.1',
      apiVersion: 'apps/v1',
      kind: 'DaemonSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Daemon Set is OK' } },
      ],
      path: {
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
        config: {
          value: 'metadata-proxy-v0.1',
          ref:
            '/overview/namespace/kube-system/workloads/daemon-sets/metadata-proxy-v0.1',
        },
      },
      hasChildren: false,
    },
  },
};

export const REAL_DATA_DEPLOYMENT: NodeDataDef = {
  edges: {
    'f604c1fe-38e2-4d55-bcc9-619b60fd0213': [
      { node: '9a44169b-ed97-461f-9c5e-1bac85deed13', edge: 'explicit' },
    ],
    '9a44169b-ed97-461f-9c5e-1bac85deed13': [
      { node: 'echo1-56b7744b6c pods', edge: 'explicit' },
    ],
    'echo1-56b7744b6c pods': [
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
      { node: 'e2fdc658-2322-40b7-907f-ba9e10df5840', edge: 'explicit' },
    ],
  },
  nodes: {
    '3c81e771-d723-403d-a19b-be7ce87ff7f2': {
      name: 'default-token-4dln7',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'default-token-4dln7',
          ref:
            '/overview/namespace/milan/config-and-storage/secrets/default-token-4dln7',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '9a44169b-ed97-461f-9c5e-1bac85deed13': {
      name: 'echo1-56b7744b6c',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'echo1-56b7744b6c',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/echo1-56b7744b6c',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'a4e5517e-0563-4158-88d3-a0492fe18cd5': {
      name: 'default',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'default',
          ref:
            '/overview/namespace/milan/config-and-storage/service-accounts/default',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'e2fdc658-2322-40b7-907f-ba9e10df5840': {
      name: 'echo1',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'echo1',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/echo1',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'echo1-56b7744b6c pods': {
      name: 'echo1-56b7744b6c pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'echo1-56b7744b6c-cwnbf': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'echo1-56b7744b6c-tdqfn': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'f604c1fe-38e2-4d55-bcc9-619b60fd0213': {
      name: 'echo1',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'echo1',
          ref: '/overview/namespace/milan/workloads/deployments/echo1',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
  },
};

export const REAL_DATA_TWO_REPLICAS: NodeDataDef = {
  edges: {
    '04ddee7a-342c-46b0-8c57-ec8682aff2ef': [
      { node: '5b287e6a-94f2-4ac3-8241-17fd87d3a114', edge: 'explicit' },
    ],
    '1b45cdc5-756b-4a0d-b8cd-54520781f0dc': [
      { node: '5b287e6a-94f2-4ac3-8241-17fd87d3a114', edge: 'explicit' },
    ],
    'elasticsearch-569cc48595 pods': [
      { node: '1b45cdc5-756b-4a0d-b8cd-54520781f0dc', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '96dcb87c-0d5e-49f8-a084-cf79e054a4bd', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'elasticsearch-dbf4fc4df pods': [
      { node: '04ddee7a-342c-46b0-8c57-ec8682aff2ef', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '96dcb87c-0d5e-49f8-a084-cf79e054a4bd', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
  },
  nodes: {
    '04ddee7a-342c-46b0-8c57-ec8682aff2ef': {
      name: 'elasticsearch-dbf4fc4df',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'elasticsearch-dbf4fc4df',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/elasticsearch-dbf4fc4df',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '1b45cdc5-756b-4a0d-b8cd-54520781f0dc': {
      name: 'elasticsearch-569cc48595',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Expected 1 replicas, but 0 are available' },
        },
      ],
      path: {
        config: {
          value: 'elasticsearch-569cc48595',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/elasticsearch-569cc48595',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '3c81e771-d723-403d-a19b-be7ce87ff7f2': {
      name: 'default-token-4dln7',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'default-token-4dln7',
          ref:
            '/overview/namespace/milan/config-and-storage/secrets/default-token-4dln7',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '5b287e6a-94f2-4ac3-8241-17fd87d3a114': {
      name: 'elasticsearch',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Expected 2 replicas, but 1 are available' },
        },
      ],
      path: {
        config: {
          value: 'elasticsearch',
          ref: '/overview/namespace/milan/workloads/deployments/elasticsearch',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '96dcb87c-0d5e-49f8-a084-cf79e054a4bd': {
      name: 'elasticsearch',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'elasticsearch',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/elasticsearch',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'a4e5517e-0563-4158-88d3-a0492fe18cd5': {
      name: 'default',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'default',
          ref:
            '/overview/namespace/milan/config-and-storage/service-accounts/default',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'elasticsearch-569cc48595 pods': {
      name: 'elasticsearch-569cc48595 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'warning',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'elasticsearch-569cc48595-s52bl': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'warning',
              },
            },
          },
        },
      ],
    },
    'elasticsearch-dbf4fc4df pods': {
      name: 'elasticsearch-dbf4fc4df pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'elasticsearch-dbf4fc4df-vnc7f': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
  },
};

export const REAL_DATA_JOB: NodeDataDef = {
  edges: {
    'contour-certgen-v1.12.0 pods': [
      { node: '3513f91b-d923-4e02-9fc6-a354666656f6', edge: 'explicit' },
      { node: '77e47f54-1599-4039-bb7b-6e638b69aeb5', edge: 'explicit' },
      { node: '8ecf1647-c947-4674-a6f4-691815f85d57', edge: 'explicit' },
    ],
  },
  nodes: {
    '3513f91b-d923-4e02-9fc6-a354666656f6': {
      name: 'contour-certgen-v1.12.0',
      apiVersion: 'batch/v1',
      kind: 'Job',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Job has succeeded 1 time' },
        },
        {
          metadata: { type: 'text' },
          config: { value: 'Job completed in 4s' },
        },
      ],
      path: {
        config: {
          value: 'contour-certgen-v1.12.0',
          ref:
            '/overview/namespace/contour-internal/workloads/jobs/contour-certgen-v1.12.0',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '77e47f54-1599-4039-bb7b-6e638b69aeb5': {
      name: 'contour-certgen-token-5ggj6',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'contour-certgen-token-5ggj6',
          ref:
            '/overview/namespace/contour-internal/config-and-storage/secrets/contour-certgen-token-5ggj6',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '8ecf1647-c947-4674-a6f4-691815f85d57': {
      name: 'contour-certgen',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'contour-certgen',
          ref:
            '/overview/namespace/contour-internal/config-and-storage/service-accounts/contour-certgen',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'contour-certgen-v1.12.0 pods': {
      name: 'contour-certgen-v1.12.0 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'warning',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'contour-certgen-v1.12.0-58hp9': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'warning',
              },
            },
          },
        },
      ],
    },
  },
};

export const REAL_DATA_INGRESS: NodeDataDef = {
  edges: {
    '2ed850cd-8bcd-4c17-9bed-45cabc70ebdf': [
      { node: '1965b0bb-86e2-4180-9746-228e5939cce2', edge: 'explicit' },
    ],
    'fd5d811f-205a-4196-a872-80c0308389c4': [
      { node: '6cf421e2-c3c8-4a88-81ef-e6d336fdb748', edge: 'explicit' },
      { node: '89b930ac-01df-4e7b-8cf9-4880ec64e887', edge: 'explicit' },
    ],
    'web-79d88c97d6 pods': [
      { node: '2ed850cd-8bcd-4c17-9bed-45cabc70ebdf', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '6cf421e2-c3c8-4a88-81ef-e6d336fdb748', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
  },
  nodes: {
    '1965b0bb-86e2-4180-9746-228e5939cce2': {
      name: 'web',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'web',
          ref: '/overview/namespace/milan/workloads/deployments/web',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '2ed850cd-8bcd-4c17-9bed-45cabc70ebdf': {
      name: 'web-79d88c97d6',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'web-79d88c97d6',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/web-79d88c97d6',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '3c81e771-d723-403d-a19b-be7ce87ff7f2': {
      name: 'default-token-4dln7',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'default-token-4dln7',
          ref:
            '/overview/namespace/milan/config-and-storage/secrets/default-token-4dln7',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '6cf421e2-c3c8-4a88-81ef-e6d336fdb748': {
      name: 'web',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'web',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/web',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '89b930ac-01df-4e7b-8cf9-4880ec64e887': {
      name: 'web2',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Service has no endpoint addresses' },
        },
      ],
      path: {
        config: {
          value: 'web2',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/web2',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'a4e5517e-0563-4158-88d3-a0492fe18cd5': {
      name: 'default',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'default',
          ref:
            '/overview/namespace/milan/config-and-storage/service-accounts/default',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'fd5d811f-205a-4196-a872-80c0308389c4': {
      name: 'example-ingress',
      apiVersion: 'networking.k8s.io/v1',
      kind: 'Ingress',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Ingress is OK' } },
      ],
      path: {
        config: {
          value: 'example-ingress',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/ingresses/example-ingress',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'web-79d88c97d6 pods': {
      name: 'web-79d88c97d6 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'web-79d88c97d6-5t5fn': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'web-79d88c97d6-7zwwf': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
  },
};

export const REAL_DATA_CRDS: NodeDataDef = {
  edges: {
    '4d2efbd7-bb91-48e0-974c-dabc5e114b13': [
      { node: '149f12bc-d4a8-4b4b-b5d5-766675fedce8', edge: 'explicit' },
      { node: 'b24a4599-477d-4b07-a992-2e622a3de8e6', edge: 'explicit' },
      { node: 'c929647a-1044-4251-81ab-147de9c0c80d', edge: 'explicit' },
    ],
    'b13a5105-4ac9-413b-92fc-e700f65bd8ce': [
      { node: '27461d6a-bd51-4a74-a1e1-00e1bb779c54', edge: 'explicit' },
      { node: 'c929647a-1044-4251-81ab-147de9c0c80d', edge: 'explicit' },
      { node: 'fc08ad66-9b99-4720-8a9b-079cec6ec17c', edge: 'explicit' },
    ],
    'b84dbc71-0bc2-4a81-8d90-8d8b3426870a': [
      { node: 'c929647a-1044-4251-81ab-147de9c0c80d', edge: 'explicit' },
      { node: 'daac8ce7-9d28-4d95-98d8-4b6a097d5ae6', edge: 'explicit' },
    ],
    'f3d53d61-f1b8-45a1-afeb-f3f18f5e2ef5': [
      { node: '9fa0495f-be04-401f-a61d-f9afd9a66d76', edge: 'explicit' },
      { node: 'c929647a-1044-4251-81ab-147de9c0c80d', edge: 'explicit' },
      { node: 'e1065101-8e56-408a-92d8-39b305d76a19', edge: 'explicit' },
    ],
  },
  nodes: {
    '149f12bc-d4a8-4b4b-b5d5-766675fedce8': {
      name: 'capi-quickstart-md-0-54w8c',
      apiVersion: 'infrastructure.cluster.x-k8s.io/v1alpha3',
      kind: 'GCPMachine',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'infrastructure.cluster.x-k8s.io/v1alpha3 GCPMachine is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-54w8c',
          ref:
            '/overview/namespace/milan/custom-resources/gcpmachines.infrastructure.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-54w8c',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '27461d6a-bd51-4a74-a1e1-00e1bb779c54': {
      name: 'capi-quickstart-md-0-hswzw',
      apiVersion: 'infrastructure.cluster.x-k8s.io/v1alpha3',
      kind: 'GCPMachine',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'infrastructure.cluster.x-k8s.io/v1alpha3 GCPMachine is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-hswzw',
          ref:
            '/overview/namespace/milan/custom-resources/gcpmachines.infrastructure.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-hswzw',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '4d2efbd7-bb91-48e0-974c-dabc5e114b13': {
      name: 'capi-quickstart-md-0-996bb7685-kgwqk',
      apiVersion: 'cluster.x-k8s.io/v1alpha3',
      kind: 'Machine',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'cluster.x-k8s.io/v1alpha3 Machine is OK' },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-996bb7685-kgwqk',
          ref:
            '/overview/namespace/milan/custom-resources/machines.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-996bb7685-kgwqk',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '9fa0495f-be04-401f-a61d-f9afd9a66d76': {
      name: 'capi-quickstart-md-0-pc6tv',
      apiVersion: 'infrastructure.cluster.x-k8s.io/v1alpha3',
      kind: 'GCPMachine',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'infrastructure.cluster.x-k8s.io/v1alpha3 GCPMachine is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-pc6tv',
          ref:
            '/overview/namespace/milan/custom-resources/gcpmachines.infrastructure.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-pc6tv',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'b13a5105-4ac9-413b-92fc-e700f65bd8ce': {
      name: 'capi-quickstart-md-0-996bb7685-lq985',
      apiVersion: 'cluster.x-k8s.io/v1alpha3',
      kind: 'Machine',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'cluster.x-k8s.io/v1alpha3 Machine is OK' },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-996bb7685-lq985',
          ref:
            '/overview/namespace/milan/custom-resources/machines.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-996bb7685-lq985',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'b24a4599-477d-4b07-a992-2e622a3de8e6': {
      name: 'capi-quickstart-md-0-vvklf',
      apiVersion: 'bootstrap.cluster.x-k8s.io/v1alpha3',
      kind: 'KubeadmConfig',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'bootstrap.cluster.x-k8s.io/v1alpha3 KubeadmConfig is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-vvklf',
          ref:
            '/overview/namespace/milan/custom-resources/kubeadmconfigs.bootstrap.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-vvklf',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'b84dbc71-0bc2-4a81-8d90-8d8b3426870a': {
      name: 'capi-quickstart-md-0',
      apiVersion: 'cluster.x-k8s.io/v1alpha3',
      kind: 'MachineDeployment',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'cluster.x-k8s.io/v1alpha3 MachineDeployment is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0',
          ref:
            '/overview/namespace/milan/custom-resources/machinedeployments.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'c929647a-1044-4251-81ab-147de9c0c80d': {
      name: 'capi-quickstart-md-0-996bb7685',
      apiVersion: 'cluster.x-k8s.io/v1alpha3',
      kind: 'MachineSet',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'cluster.x-k8s.io/v1alpha3 MachineSet is OK' },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-996bb7685',
          ref:
            '/overview/namespace/milan/custom-resources/machinesets.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-996bb7685',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'daac8ce7-9d28-4d95-98d8-4b6a097d5ae6': {
      name: 'capi-quickstart',
      apiVersion: 'cluster.x-k8s.io/v1alpha3',
      kind: 'Cluster',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'cluster.x-k8s.io/v1alpha3 Cluster is OK' },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart',
          ref:
            '/overview/namespace/milan/custom-resources/clusters.cluster.x-k8s.io/v1alpha3/capi-quickstart',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'e1065101-8e56-408a-92d8-39b305d76a19': {
      name: 'capi-quickstart-md-0-4rqgc',
      apiVersion: 'bootstrap.cluster.x-k8s.io/v1alpha3',
      kind: 'KubeadmConfig',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'bootstrap.cluster.x-k8s.io/v1alpha3 KubeadmConfig is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-4rqgc',
          ref:
            '/overview/namespace/milan/custom-resources/kubeadmconfigs.bootstrap.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-4rqgc',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'f3d53d61-f1b8-45a1-afeb-f3f18f5e2ef5': {
      name: 'capi-quickstart-md-0-996bb7685-wh2tx',
      apiVersion: 'cluster.x-k8s.io/v1alpha3',
      kind: 'Machine',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'cluster.x-k8s.io/v1alpha3 Machine is OK' },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-996bb7685-wh2tx',
          ref:
            '/overview/namespace/milan/custom-resources/machines.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-996bb7685-wh2tx',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'fc08ad66-9b99-4720-8a9b-079cec6ec17c': {
      name: 'capi-quickstart-md-0-8rwp6',
      apiVersion: 'bootstrap.cluster.x-k8s.io/v1alpha3',
      kind: 'KubeadmConfig',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: {
            value: 'bootstrap.cluster.x-k8s.io/v1alpha3 KubeadmConfig is OK',
          },
        },
      ],
      path: {
        config: {
          value: 'capi-quickstart-md-0-8rwp6',
          ref:
            '/overview/namespace/milan/custom-resources/kubeadmconfigs.bootstrap.cluster.x-k8s.io/v1alpha3/capi-quickstart-md-0-8rwp6',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
  },
};

export const REAL_DATA_CRDS2: NodeDataDef = {
  edges: {
    '039fb462-0c29-47b3-81aa-dcf0b7ab695e': [
      { node: 'loader-69fb98c8b5 pods', edge: 'explicit' },
    ],
    '04ddee7a-342c-46b0-8c57-ec8682aff2ef': [
      { node: '5b287e6a-94f2-4ac3-8241-17fd87d3a114', edge: 'explicit' },
    ],
    '11333700-3e1a-4633-b214-f5c579b02572': [
      { node: '20b000c0-64c3-4935-ad98-32c4f001fdf8', edge: 'explicit' },
    ],
    '1b45cdc5-756b-4a0d-b8cd-54520781f0dc': [
      { node: '5b287e6a-94f2-4ac3-8241-17fd87d3a114', edge: 'explicit' },
    ],
    '2ed850cd-8bcd-4c17-9bed-45cabc70ebdf': [
      { node: '1965b0bb-86e2-4180-9746-228e5939cce2', edge: 'explicit' },
    ],
    '37cc6a74-8437-496b-9c4c-932433a1ce48': [
      { node: '5518f305-d8e4-4a54-a637-ef131564a38a', edge: 'explicit' },
    ],
    '638619e3-5261-4c54-b32e-83cdca4bd76c': [
      { node: 'ca762c1b-05da-48f5-99e6-8f9ce2aa748f', edge: 'explicit' },
    ],
    '71ca4e15-566a-4cbc-930e-102e5cfd0f4e': [
      { node: '31a54305-1c35-40ab-97a9-a27a8c102c0d', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
    ],
    '73a42d8e-ddba-4ca4-8db1-02e41ded0e08': [
      { node: 'jobposting-57bd4c8596 pods', edge: 'explicit' },
    ],
    '96dcb87c-0d5e-49f8-a084-cf79e054a4bd': [
      { node: 'elasticsearch-569cc48595 pods', edge: 'explicit' },
    ],
    '9a44169b-ed97-461f-9c5e-1bac85deed13': [
      { node: 'f604c1fe-38e2-4d55-bcc9-619b60fd0213', edge: 'explicit' },
    ],
    'ae87be09-f766-4011-b986-d9d56e5224a7': [
      { node: 'ea66893f-5967-4a7b-a5d9-bdc975c29ab1', edge: 'explicit' },
    ],
    'c0d6ad21-bfb5-4745-9345-a685b1e70919': [
      { node: '0da1180b-b908-42fb-b375-ed3bddc5e13a', edge: 'explicit' },
    ],
    'crawler-67cf8bbcdd pods': [
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
      { node: 'f927d239-7b0a-4866-ac71-245a7c91bcc8', edge: 'explicit' },
    ],
    'e2fdc658-2322-40b7-907f-ba9e10df5840': [
      { node: 'echo1-56b7744b6c pods', edge: 'explicit' },
    ],
    'echo1-56b7744b6c pods': [
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: '9a44169b-ed97-461f-9c5e-1bac85deed13', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'elasticsearch-569cc48595 pods': [
      { node: '1b45cdc5-756b-4a0d-b8cd-54520781f0dc', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'elasticsearch-dbf4fc4df pods': [
      { node: '04ddee7a-342c-46b0-8c57-ec8682aff2ef', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: '96dcb87c-0d5e-49f8-a084-cf79e054a4bd', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'f927d239-7b0a-4866-ac71-245a7c91bcc8': [
      { node: 'da2f1f77-75a0-4d56-800a-3d5118ccd446', edge: 'explicit' },
    ],
    'fd5d811f-205a-4196-a872-80c0308389c4': [
      { node: '6cf421e2-c3c8-4a88-81ef-e6d336fdb748', edge: 'explicit' },
      { node: '89b930ac-01df-4e7b-8cf9-4880ec64e887', edge: 'explicit' },
    ],
    'jobposting-57bd4c8596 pods': [
      { node: '11333700-3e1a-4633-b214-f5c579b02572', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'kafka pods': [
      { node: '0bf159aa-01ea-4742-a6a2-becef1178827', edge: 'explicit' },
      { node: '14eda8ed-87c3-4aa1-a3cb-9f4279704fc5', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: '9cd5e4f0-4979-4119-9c93-7df18bd88059', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'loader-69fb98c8b5 pods': [
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
      { node: 'ae87be09-f766-4011-b986-d9d56e5224a7', edge: 'explicit' },
    ],
    'recruiter-54f94f7b87 pods': [
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '638619e3-5261-4c54-b32e-83cdca4bd76c', edge: 'explicit' },
      { node: '74357386-8506-4393-a1c1-79895cb1fa21', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'rolling-test-574b87c764 pods': [
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
      { node: 'c0d6ad21-bfb5-4745-9345-a685b1e70919', edge: 'explicit' },
    ],
    'web-79d88c97d6 pods': [
      { node: '2ed850cd-8bcd-4c17-9bed-45cabc70ebdf', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '6cf421e2-c3c8-4a88-81ef-e6d336fdb748', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
    'zookeeper-66b5f99f97 pods': [
      { node: '37cc6a74-8437-496b-9c4c-932433a1ce48', edge: 'explicit' },
      { node: '3c81e771-d723-403d-a19b-be7ce87ff7f2', edge: 'explicit' },
      { node: '77bf3dd6-9dc3-49b7-873b-911b21c23bb8', edge: 'explicit' },
      { node: '8ac76613-f4be-4b3c-84fe-d6f36467d2f6', edge: 'explicit' },
      { node: 'a4e5517e-0563-4158-88d3-a0492fe18cd5', edge: 'explicit' },
    ],
  },
  nodes: {
    '039fb462-0c29-47b3-81aa-dcf0b7ab695e': {
      name: 'loader',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'loader',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/loader',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '04ddee7a-342c-46b0-8c57-ec8682aff2ef': {
      name: 'elasticsearch-dbf4fc4df',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'elasticsearch-dbf4fc4df',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/elasticsearch-dbf4fc4df',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '0bf159aa-01ea-4742-a6a2-becef1178827': {
      name: 'kafka-config',
      apiVersion: 'v1',
      kind: 'ConfigMap',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 ConfigMap is OK' } },
      ],
      path: {
        config: {
          value: 'kafka-config',
          ref:
            '/overview/namespace/milan/config-and-storage/config-maps/kafka-config',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '0da1180b-b908-42fb-b375-ed3bddc5e13a': {
      name: 'rolling-test',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'rolling-test',
          ref: '/overview/namespace/milan/workloads/deployments/rolling-test',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '11333700-3e1a-4633-b214-f5c579b02572': {
      name: 'jobposting-57bd4c8596',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'jobposting-57bd4c8596',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/jobposting-57bd4c8596',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '14eda8ed-87c3-4aa1-a3cb-9f4279704fc5': {
      name: 'kafka-headless',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'kafka-headless',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/kafka-headless',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '1965b0bb-86e2-4180-9746-228e5939cce2': {
      name: 'web',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'web',
          ref: '/overview/namespace/milan/workloads/deployments/web',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '1b45cdc5-756b-4a0d-b8cd-54520781f0dc': {
      name: 'elasticsearch-569cc48595',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Expected 1 replicas, but 0 are available' },
        },
      ],
      path: {
        config: {
          value: 'elasticsearch-569cc48595',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/elasticsearch-569cc48595',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '20b000c0-64c3-4935-ad98-32c4f001fdf8': {
      name: 'jobposting',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'jobposting',
          ref: '/overview/namespace/milan/workloads/deployments/jobposting',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '2ed850cd-8bcd-4c17-9bed-45cabc70ebdf': {
      name: 'web-79d88c97d6',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'web-79d88c97d6',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/web-79d88c97d6',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '31a54305-1c35-40ab-97a9-a27a8c102c0d': {
      name: 'default',
      apiVersion: 'eventing.knative.dev/v1',
      kind: 'Broker',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Broker is being deleted' },
        },
      ],
      path: {
        config: {
          value: 'default',
          ref:
            '/overview/namespace/milan/custom-resources/brokers.eventing.knative.dev/v1/default',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '37cc6a74-8437-496b-9c4c-932433a1ce48': {
      name: 'zookeeper-66b5f99f97',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'zookeeper-66b5f99f97',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/zookeeper-66b5f99f97',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '3c81e771-d723-403d-a19b-be7ce87ff7f2': {
      name: 'default-token-4dln7',
      apiVersion: 'v1',
      kind: 'Secret',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'v1 Secret is OK' } },
      ],
      path: {
        config: {
          value: 'default-token-4dln7',
          ref:
            '/overview/namespace/milan/config-and-storage/secrets/default-token-4dln7',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '5518f305-d8e4-4a54-a637-ef131564a38a': {
      name: 'zookeeper',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'zookeeper',
          ref: '/overview/namespace/milan/workloads/deployments/zookeeper',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '5b287e6a-94f2-4ac3-8241-17fd87d3a114': {
      name: 'elasticsearch',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Expected 2 replicas, but 1 are available' },
        },
      ],
      path: {
        config: {
          value: 'elasticsearch',
          ref: '/overview/namespace/milan/workloads/deployments/elasticsearch',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '638619e3-5261-4c54-b32e-83cdca4bd76c': {
      name: 'recruiter-54f94f7b87',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'recruiter-54f94f7b87',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/recruiter-54f94f7b87',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '6cf421e2-c3c8-4a88-81ef-e6d336fdb748': {
      name: 'web',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'web',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/web',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '71ca4e15-566a-4cbc-930e-102e5cfd0f4e': {
      name: 'default-kne-trigger',
      apiVersion: 'messaging.knative.dev/v1',
      kind: 'InMemoryChannel',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'messaging.knative.dev/v1 InMemoryChannel is OK' },
        },
      ],
      path: {
        config: {
          value: 'default-kne-trigger',
          ref:
            '/overview/namespace/milan/custom-resources/inmemorychannels.messaging.knative.dev/v1/default-kne-trigger',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '73a42d8e-ddba-4ca4-8db1-02e41ded0e08': {
      name: 'jobposting',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'jobposting',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/jobposting',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '74357386-8506-4393-a1c1-79895cb1fa21': {
      name: 'recruiter',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'recruiter',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/recruiter',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '77bf3dd6-9dc3-49b7-873b-911b21c23bb8': {
      name: 'zk-headless',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'zk-headless',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/zk-headless',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '89b930ac-01df-4e7b-8cf9-4880ec64e887': {
      name: 'web2',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Service has no endpoint addresses' },
        },
      ],
      path: {
        config: {
          value: 'web2',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/web2',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '8ac76613-f4be-4b3c-84fe-d6f36467d2f6': {
      name: 'default-kne-trigger-kn-channel',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'default-kne-trigger-kn-channel',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/default-kne-trigger-kn-channel',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '96dcb87c-0d5e-49f8-a084-cf79e054a4bd': {
      name: 'elasticsearch',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'elasticsearch',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/elasticsearch',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '9a44169b-ed97-461f-9c5e-1bac85deed13': {
      name: 'echo1-56b7744b6c',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'echo1-56b7744b6c',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/echo1-56b7744b6c',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    '9cd5e4f0-4979-4119-9c93-7df18bd88059': {
      name: 'kafka',
      apiVersion: 'apps/v1',
      kind: 'StatefulSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Stateful Set is OK' } },
      ],
      path: {
        config: {
          value: 'kafka',
          ref: '/overview/namespace/milan/workloads/stateful-sets/kafka',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'a4e5517e-0563-4158-88d3-a0492fe18cd5': {
      name: 'default',
      apiVersion: 'v1',
      kind: 'ServiceAccount',
      status: 'ok',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'v1 ServiceAccount is OK' },
        },
      ],
      path: {
        config: {
          value: 'default',
          ref:
            '/overview/namespace/milan/config-and-storage/service-accounts/default',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'ae87be09-f766-4011-b986-d9d56e5224a7': {
      name: 'loader-69fb98c8b5',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'loader-69fb98c8b5',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/loader-69fb98c8b5',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'c0d6ad21-bfb5-4745-9345-a685b1e70919': {
      name: 'rolling-test-574b87c764',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Replica Set is OK' } },
      ],
      path: {
        config: {
          value: 'rolling-test-574b87c764',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/rolling-test-574b87c764',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'ca762c1b-05da-48f5-99e6-8f9ce2aa748f': {
      name: 'recruiter',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'recruiter',
          ref: '/overview/namespace/milan/workloads/deployments/recruiter',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'crawler-67cf8bbcdd pods': {
      name: 'crawler-67cf8bbcdd pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'crawler-67cf8bbcdd-2hltp': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'da2f1f77-75a0-4d56-800a-3d5118ccd446': {
      name: 'crawler',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'error',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'No replicas exist for this deployment' },
        },
      ],
      path: {
        config: {
          value: 'crawler',
          ref: '/overview/namespace/milan/workloads/deployments/crawler',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'e2fdc658-2322-40b7-907f-ba9e10df5840': {
      name: 'echo1',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Service is OK' } },
      ],
      path: {
        config: {
          value: 'echo1',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/services/echo1',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'ea66893f-5967-4a7b-a5d9-bdc975c29ab1': {
      name: 'loader',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'loader',
          ref: '/overview/namespace/milan/workloads/deployments/loader',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'echo1-56b7744b6c pods': {
      name: 'echo1-56b7744b6c pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'echo1-56b7744b6c-cwnbf': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'echo1-56b7744b6c-tdqfn': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'elasticsearch-569cc48595 pods': {
      name: 'elasticsearch-569cc48595 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'warning',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'elasticsearch-569cc48595-s52bl': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'warning',
              },
            },
          },
        },
      ],
    },
    'elasticsearch-dbf4fc4df pods': {
      name: 'elasticsearch-dbf4fc4df pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'elasticsearch-dbf4fc4df-vnc7f': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'f604c1fe-38e2-4d55-bcc9-619b60fd0213': {
      name: 'echo1',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Deployment is OK' } },
      ],
      path: {
        config: {
          value: 'echo1',
          ref: '/overview/namespace/milan/workloads/deployments/echo1',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'f927d239-7b0a-4866-ac71-245a7c91bcc8': {
      name: 'crawler-67cf8bbcdd',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'warning',
      details: [
        {
          metadata: { type: 'text' },
          config: { value: 'Expected 1 replicas, but 0 are available' },
        },
      ],
      path: {
        config: {
          value: 'crawler-67cf8bbcdd',
          ref:
            '/overview/namespace/milan/workloads/replica-sets/crawler-67cf8bbcdd',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'fd5d811f-205a-4196-a872-80c0308389c4': {
      name: 'example-ingress',
      apiVersion: 'networking.k8s.io/v1',
      kind: 'Ingress',
      status: 'ok',
      details: [
        { metadata: { type: 'text' }, config: { value: 'Ingress is OK' } },
      ],
      path: {
        config: {
          value: 'example-ingress',
          ref:
            '/overview/namespace/milan/discovery-and-load-balancing/ingresses/example-ingress',
        },
        metadata: {
          type: 'link',
          title: [{ metadata: { type: 'text' }, config: { value: '' } }],
        },
      },
    },
    'jobposting-57bd4c8596 pods': {
      name: 'jobposting-57bd4c8596 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'jobposting-57bd4c8596-pvtd4': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'kafka pods': {
      name: 'kafka pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'kafka-0': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'loader-69fb98c8b5 pods': {
      name: 'loader-69fb98c8b5 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'loader-69fb98c8b5-fzsd7': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'recruiter-54f94f7b87 pods': {
      name: 'recruiter-54f94f7b87 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'recruiter-54f94f7b87-t6vh7': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'rolling-test-574b87c764 pods': {
      name: 'rolling-test-574b87c764 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'rolling-test-574b87c764-dw6cj': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'rolling-test-574b87c764-jwljk': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'rolling-test-574b87c764-rtfjb': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'web-79d88c97d6 pods': {
      name: 'web-79d88c97d6 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'web-79d88c97d6-5t5fn': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
              'web-79d88c97d6-7zwwf': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
    'zookeeper-66b5f99f97 pods': {
      name: 'zookeeper-66b5f99f97 pods',
      apiVersion: 'v1',
      kind: 'Pod',
      status: 'ok',
      details: [
        {
          metadata: { type: 'podStatus' },
          config: {
            pods: {
              'zookeeper-66b5f99f97-jqg9h': {
                details: [
                  { metadata: { type: 'text' }, config: { value: '' } },
                ],
                status: 'ok',
              },
            },
          },
        },
      ],
    },
  },
};
