import NodeShape = cytoscape.Css.NodeShape;

export const ELEMENTS_STYLE = [
  {
    selector: 'node',
    css: {
      shape: node => nodeShape(node),
      width: 'data(width)',
      height: 'data(height)',
      content: 'data(label)',
      'background-color': '#F2F2F2',
      color: 'black',
      'border-color': 'black',
      'border-width': '2px',
      'border-style': 'solid',
      fontSize: 16,
      ghost: 'no',
      'text-wrap': 'wrap',
      'text-valign': 'top',
      'text-halign': 'center',
      'text-margin-y': 20,
      padding: '10px',
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
    selector: '.unbundled',
    css: {
      'curve-style': 'unbundled-bezier',
      'source-endpoint': '90deg',
      'target-endpoint': '270deg',
    },
  },
  {
    selector: '.pod',
    css: {
      ghost: 'yes',
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
      visibility: 'hidden',
    },
  },
];

function nodeShape(node): NodeShape {
  return node.data('shape');
}
