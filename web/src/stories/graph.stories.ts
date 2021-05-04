import { ELEMENTS_STYLE } from './graph.data';
import {
  BaseShape,
  Deployment,
  Edge,
  Pod,
  Port,
  ReplicaSet,
  Secret,
  Service,
  ServiceAccount,
} from '../app/modules/shared/components/presentation/cytoscape2/shape';
import { Meta, Story } from '@storybook/angular/types-6-0';

const layout = {
  name: 'cose-bilkent',
  padding: 30,
  fit: false,
  animateFilter: () => false,
};

const zoom = {
  min: 0.1,
  max: 2.0,
};

const shapes: BaseShape[] = [
  new Deployment('glyph0', 'Deployment', true),
  new Secret('glyph2', 'Secret', true),
  new ServiceAccount('glyph3', 'ServiceAccount', false),
  new Service('glyph1', 'Service', true),
  new ReplicaSet('glyph10', 'ReplicaSet: 3', true, 'glyph0'),
  new Pod('glyph30', 'Pods', true, 'glyph10'),
  new Port('glyph20', 'image: nginx', 'left', 'port', 'glyph0'),
  new Port('glyph21', 'metadata.annotations', 'right', 'port', 'glyph2'),
  new Port('glyph41', 'app: demo', 'left', 'label', 'glyph30'),
  new Port('glyph42', 'app: demo', 'right', 'selector', 'glyph1'),
  new Port('glyph50', 'name', 'right', 'port', 'glyph3'),
  new Port('glyph51', 'serviceAccount', 'left', 'port', 'glyph30'),
  new Port('glyph52', 'secrets.name', 'left', 'port', 'glyph3'),
  new Edge('glyph42-glyph41', 'glyph42', 'glyph41'),
  new Edge('glyph52-glyph21', 'glyph52', 'glyph21'),
  new Edge('glyph50-glyph51', 'glyph50', 'glyph51', 'unbundled'),
  // new Edge('glyph1-glyph30', 'glyph1', 'glyph30'),
  // new Edge('glyph3-glyph2', 'glyph3', 'glyph2'),
  // new Edge('glyph3-glyph30', 'glyph3', 'glyph30'),
];
const style = ELEMENTS_STYLE;

export const graphStory: Story = args => {
  return {
    props: {
      elements: args.elements,
      layout,
      zoom,
      style,
    },
    template: `
      <div class="main-container">
          <div class="content-container">
              <div class="content-area" style="background-color: white;">
                  <app-cytoscape2
                    [elements]="elements"
                    [layout]="layout"
                    [zoom]="zoom"
                    [style]="style">
                  </app-cytoscape2>
              </div>
          </div>
      </div>
      `,
  };
};

graphStory.storyName = 'Resource View prototype';

graphStory.argTypes = {
  elements: {
    control: {
      type: 'object',
    },
  },
};

graphStory.args = {
  elements: shapes.map(shape => shape.toNode(shapes)),
};
