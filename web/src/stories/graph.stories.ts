import { storiesOf } from '@storybook/angular';
import {ELEMENTS_DATA, ELEMENTS_STYLE} from "./graph.data";
import octant from '../app/modules/shared/components/presentation/cytoscape/octant.layout'
import cytoscape from 'cytoscape';
import coseBilkent from 'cytoscape-cose-bilkent';

cytoscape.use( coseBilkent );
cytoscape( 'layout', 'octant', octant );

const layout = { name: 'cose-bilkent', padding: 30, fit: true, animate: true};

const zoom= {
    "min": 0.5,
    "max": 3
  };

const elements= ELEMENTS_DATA;
const style= ELEMENTS_STYLE;

storiesOf('Resources', module).add('Resource View', () => ({
  props: {
    elements: elements,
    layout: layout,
    zoom: zoom,
    style: style,
  },
  template: `
    <div class="main-container">
        <div class="content-container">
            <div class="content-area" style="background-color: white;">
                <app-cytoscape
                  [elements]="elements" 
                  [layout]="layout" 
                  [zoom]="zoom" 
                  [style]="style">
                </app-cytoscape>
            </div>
        </div>
    </div>
    `,
}));
