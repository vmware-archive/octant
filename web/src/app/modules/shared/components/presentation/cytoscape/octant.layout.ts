import {
  AnimatedLayoutOptions,
  BaseLayoutOptions,
  NodeSingular,
} from 'cytoscape';
import cytoscape from 'cytoscape';

export interface OctantLayoutOptions
  extends BaseLayoutOptions,
    AnimatedLayoutOptions {
  name: 'octant';
  fit: boolean;
  padding?: number;
}

const defaults = {
  positions: undefined,
  zoom: undefined,
  pan: undefined,
  fit: true,
  padding: 30,
  animate: false,
  animationDuration: 500,
  animationEasing: undefined,
  animateFilter(node, i) {
    return true;
  },
  ready: undefined,
  stop: undefined,
  transform(node, position) {
    return position;
  },
};

export function positionChildren(
  cytoscape: cytoscape.Core,
  node: cytoscape.NodeSingular
) {
  const offset = {
    x: node.position().x - node.data('x'),
    y: node.position().y - node.data('y'),
  };
  moveChildren(cytoscape, node, offset);
  moveNode(node, offset);

  const options: OctantLayoutOptions = { name: 'octant', fit: false };
  cytoscape.nodes().layout(options).run();
}

export function hideChildren(
  cytoscape: cytoscape.Core,
  node: cytoscape.NodeSingular
) {
  const children = cytoscape.nodes(`[owner = "${node.data('id')}"]`);

  children.map(child => {
    hideChildren(cytoscape, child);
    child.style('visibility', 'hidden');
  });
}

function OctantLayout(options) {
  this.options = { ...defaults, ...options };
}

const isFn = obj => obj != null && typeof obj === typeof function () {};

OctantLayout.prototype.run = function () {
  const options = this.options;
  const eles = options.eles;

  const nodes = eles.nodes();
  const posIsFn = isFn(options.positions);

  function getPosition(node) {
    if (options.positions == null) {
      return { x: node.position().x, y: node.position().y };
    }

    if (posIsFn) {
      return options.positions(node);
    }

    const pos = options.positions[node._private.data.id];

    if (pos == null) {
      return null;
    }

    return pos;
  }

  nodes.layoutPositions(this, options, node => {
    const position = getPosition(node);

    if (node.locked() || position === null) {
      return { x: 0, y: 0 };
    }

    return { x: node.data('x'), y: node.data('y') };
  });

  return this; // chaining
};

function moveChildren(
  cytoscape: cytoscape.Core,
  node: cytoscape.NodeSingular,
  offset
) {
  const children = cytoscape.nodes(`[owner = "${node.data('id')}"]`);

  children.map(child => {
    moveNode(child, offset);
    child.style('visibility', 'visible');
    moveChildren(cytoscape, child, offset);
  });
  return children;
}

function moveNode(node: NodeSingular, offset: cytoscape.Position) {
  node.data('x', node.data('x') + offset.x);
  node.data('y', node.data('y') + offset.y);
}

export default OctantLayout;
