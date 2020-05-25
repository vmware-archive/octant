import NodeShape = cytoscape.Css.NodeShape;

export const ELEMENTS_DATA = {
  "nodes": [
    {
      "data": {
        "id": "glyph0",
        "label": "Deployment",
        "width": 800,
        "height": 600,
        "x": 750,
        "y": 450,
        "hasChildren": true,
        "shape": "rectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": false,
      "classes": "deployment"
    },
    {
      "data": {
        "id": "glyph2",
        "label": "Secret",
        "width": 350,
        "height": 200,
        "x": 100,
        "y": 1000,
        "hasChildren": true,
        "shape": "roundrectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": false,
      "classes": "secret"
    },
    {
      "data": {
        "id": "glyph21",
        "label": "metadata.annotations",
        "owner": "glyph2",
        "width": "label",
        "height": "label",
        "x": 195,
        "y": 1005,
        "hasChildren": false,
        "shape": "rectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": false,
      "locked": false,
      "grabbable": false,
      "pannable": false,
      "classes": "port"
    },
    {
      "data": {
        "id": "glyph3",
        "label": "ServiceAccount",
        "width": 350,
        "height": 200,
        "x": 600,
        "y": 1000,
        "hasChildren": true,
        "shape": "roundrectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": false,
      "classes": "secret"
    },
    {
      "data": {
        "id": "glyph1",
        "label": "Service",
        "width": 350,
        "height": 200,
        "x": 0,
        "y": 400,
        "hasChildren": true,
        "shape": "roundrectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": false,
      "classes": "secret"
    },
    {
      "data": {
        "id": "glyph10",
        "label": "ReplicaSet: 3",
        "owner": "glyph0",
        "width": 600,
        "height": 400,
        "x": 800,
        "y": 500,
        "hasChildren": true,
        "shape": "rectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": false,
      "classes": "replicaset"
    },
    {
      "data": {
        "id": "glyph20",
        "label": "image: nginx",
        "owner": "glyph0",
        "width": "label",
        "height": "label",
        "x": 400,
        "y": 200,
        "hasChildren": false,
        "shape": "rectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": false,
      "locked": false,
      "grabbable": false,
      "pannable": false,
      "classes": "port"
    },
    {
      "data": {
        "id": "glyph30",
        "label": "Pod",
        "owner": "glyph10",
        "width": 350,
        "height": 200,
        "x": 875,
        "y": 525,
        "hasChildren": false,
        "shape": "roundrectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": false,
      "classes": "pod"
    },
    {
      "data": {
        "id": "glyph41",
        "label": "app: demo",
        "owner": "glyph30",
        "width": "label",
        "height": "label",
        "x": 740,
        "y": 525,
        "hasChildren": false,
        "shape": "rectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": false,
      "locked": false,
      "grabbable": false,
      "pannable": false,
      "classes": "label"
    },
    {
      "data": {
        "id": "glyph42",
        "label": "app: demo",
        "owner": "glyph1",
        "width": "label",
        "height": "label",
        "x": 135,
        "y": 400,
        "hasChildren": false,
        "shape": "rectangle",
      },
      "group": "nodes",
      "removed": false,
      "selected": false,
      "selectable": false,
      "locked": false,
      "grabbable": false,
      "pannable": false,
      "classes": "selector"
    },
  ],
  "edges": [
    {
      "data": {
        "id": "glyph1-glyph30",
        "class": "consumption",
        "source": "glyph1",
        "target": "glyph30",
      },
      "group": "edges",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": true,
      "classes": ""
    },
    {
      "data": {
        "id": "glyph3-glyph2",
        "class": "consumption",
        "source": "glyph3",
        "target": "glyph2",
      },
      "group": "edges",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": true,
      "classes": ""
    },
    {
      "data": {
        "id": "glyph3-glyph30",
        "class": "consumption",
        "source": "glyph3",
        "target": "glyph30",
      },
      "group": "edges",
      "removed": false,
      "selected": false,
      "selectable": true,
      "locked": false,
      "grabbable": true,
      "pannable": true,
      "classes": ""
    },
  ]
};

export const ELEMENTS_STYLE = [
  {
    selector: 'node',
    css: {
      shape: (node) => nodeShape(node),
      width: 'data(width)',
      height: 'data(height)',
      content: 'data(label)',
      'background-color': '#F2F2F2',
      color: 'black',
      'border-color': 'black',
      'border-width': '2px',
      'border-style': 'solid',
      "fontSize": 16,
      'ghost': 'no',
      'text-wrap': 'wrap',
      'text-valign': 'top',
      'text-halign': 'center',
      'text-margin-y': 20,
      'padding': '10px',
    },
  },

  {
    selector: 'node:selected',
    css: {
      'curve-style': 'bezier',
      'border-width': 1,
      'border-color': '#313131',
      'border-style': 'solid',
    },
  },
  {
    selector: 'edge',
    css: {
      'curve-style': 'bezier',
      opacity: 1,
      width: 1.5,
      'line-color': 'black',
      'z-compound-depth': 'top',
      'source-arrow-color': 'black',
      'source-arrow-fill': 'hollow',
      'source-arrow-shape': 'tee',
      'target-arrow-color': 'black',
      'target-arrow-fill': 'hollow',
      'target-arrow-shape': 'triangle-backcurve',
      'arrow-scale': 2,
    },
  },
  {
    selector: '.pod',
    css: {
      'ghost': 'yes',
      'ghost-opacity': 1,
      'ghost-offset-x': 10,
      'ghost-offset-y': 10,
      'font-size': 24,
      'text-margin-y': 30,
      'border-width': 1.5,
    },
  },
  {
    selector: '.deployment',
    css: {
      'font-size': 32,
      'text-margin-y': 38,
    },
  },
  {
    selector: '.secret',
    css: {
      'font-size': 24,
      'text-margin-y': 30,
    },
  },
  {
    selector: '.replicaset',
    css: {
      'font-size': 32,
      'text-margin-y': 38,
      'border-style': 'dashed',
      'border-width': 3,
    },
  },
  {
    selector: '.label',
    css: {
      'background-color': '#13C6CE',
      'border-width': '0px',
    },
  },
  {
    selector: '.port',
    css: {
      'border-width': '0px',
    },
  },
  {
    selector: '.selector',
    css: {
      'background-color': '#F9C011',
      'border-width': '0px',
    },
  },
  {
    selector: '[owner]',
    css: {
      "visibility": "hidden",
    },
  },
];

function nodeShape(node):NodeShape {
  return(node.data('shape'))
}
