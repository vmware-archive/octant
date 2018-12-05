import Grid from './grid'

describe('Grid', () => {
  describe('isSink', () => {
    const cases = [
      ['pods-rs0', true],
      ['rs0', false],
      ['rs1', false],
      ['rs2', false],
      ['rs3', false],
      ['s1', false],
      ['d1', false],
    ]

    test.each(cases)('isSink(\'%s\')', (name, expected) => {
      const grid = new Grid(complex.adjacencyList, complex.objects)
      expect(grid.isSink(name)).toEqual(expected)
    })
  })

  describe('create', () => {
    describe('with valid settings', () => {
      test('in general', () => {
        const grid = new Grid(complex.adjacencyList, complex.objects)
        const rows = grid.create()

        const expected = [
          new Set(['s1']),
          new Set(['pods-rs0']),
          new Set(['rs0', 'rs1', 'rs2', 'rs3']),
          new Set(['d1']),
        ]

        expect(rows).toEqual(expected)
      })
    })

    describe('with invalid settings', () => {
      test('in general', () => {
        const grid = new Grid(
          missingRSEdges.adjacencyList,
          missingRSEdges.objects,
        )
        const rows = grid.create()

        const expected = [
          new Set(['rs1', 'pods1']),
          new Set(['rs2']),
          new Set(['d']),
        ]

        expect(rows).toEqual(expected)
      })
    })
  })

  describe('sort', () => {
    describe('with correct data', () => {
      test('sorts correctly', () => {
        const grid = new Grid(complex.adjacencyList, complex.objects)
        const sorted = grid.sort()

        const expected = ['pods-rs0', 'rs0', 'rs1', 'rs2', 's1', 'rs3', 'd1']

        expect(sorted).toEqual(expected)
      })
    })

    describe('with missing edges', () => {
      test('sorts with issues', () => {
        const grid = new Grid(
          missingRSEdges.adjacencyList,
          missingRSEdges.objects,
        )
        const sorted = grid.sort()

        const expected = ['rs1', 'pods1', 'rs2', 'd']

        expect(sorted).toEqual(expected)
      })
    })
  })

  describe('row for node', () => {
    const cases = [
      ['rs0', [new Set(['pods-rs0'])], 1],
      ['s1', [new Set(['pods-rs0'])], -1],
      ['d1', [new Set(['pods-rs0']), new Set(['rs0'])], 2],
    ]

    test.each(cases)('returns row for name %s', (name, rows, expected) => {
      const grid = new Grid(complex.adjacencyList, complex.objects)
      const row = grid.rowForNode(rows, name)
      expect(row).toEqual(expected)
    })
  })
})

const complex = {
  type: 'resourceviewer',
  selected: 'd1',
  adjacencyList: {
    rs0: [{ node: 'pods-rs0', edge: 'explicit' }],
    rs1: [{ node: 'pods-rs0', edge: 'explicit' }],
    rs2: [{ node: 'pods-rs0', edge: 'explicit' }],
    rs3: [{ node: 'pods-rs0', edge: 'explicit' }],
    s1: [{ node: 'pods-rs0', edge: 'implicit' }],
    d1: [
      { node: 'rs0', edge: 'explicit' },
      { node: 'rs1', edge: 'explicit' },
      { node: 'rs2', edge: 'explicit' },
      { node: 'rs3', edge: 'explicit' },
    ],
  },
  objects: {
    rs0: {
      name: 'grafana-6d4fd8c49',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
    },
    rs1: {
      name: 'grafana-99c8784f6',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
    },
    rs2: {
      name: 'grafana-6b5b79d6cf',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
    },
    s1: {
      name: 'grafana',
      apiVersion: 'v1',
      kind: 'Service',
      status: 'ok',
      isNetwork: true,
    },
    d1: {
      name: 'grafana',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
    },
    rs3: {
      name: 'grafana-d69f77cc4',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
    },
    'pods-rs0': {
      name: 'pods-rs0',
      apiVersion: 'v1',
      kind: 'pods',
      status: 'ok',
    },
  },
}

const missingRSEdges = {
  type: 'resourceviewer',
  selected: 'd1',
  adjacencyList: {
    d: [{ node: 'rs1', edge: 'explicit' }, { node: 'rs2', edge: 'explicit' }],
    rs2: [{ node: 'pods1', edge: 'explicit' }],
  },
  objects: {
    d: {
      name: 'nginx-apps-v1',
      apiVersion: 'apps/v1',
      kind: 'Deployment',
      status: 'ok',
    },
    rs1: {
      name: 'nginx-apps-v1-86d59dd769',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
    },
    rs2: {
      name: 'nginx-apps-v1-c97df9bdd',
      apiVersion: 'apps/v1',
      kind: 'ReplicaSet',
      status: 'ok',
    },
    pods1: {
      name: 'pods1',
      apiVersion: 'v1',
      kind: 'pods',
      status: 'ok',
    },
  },
}
