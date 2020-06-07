import { storiesOf } from '@storybook/angular';
import {number} from '@storybook/addon-knobs';
import {QuadrantComponent} from "../app/modules/shared/components/presentation/quadrant/quadrant.component";


storiesOf('Quadrant', module).add('Quadrant component', () => ({
  props: {
    view: {
      metadata: {
        type: "quadrant",
        title: [{metadata: {type: "text"}, config: {value: "Status"}}]
      },
      config: {
        nw: {value: number('NW value', 1), label: "Running"},
        ne: {value: number('NE value', 0), label: "Waiting"},
        se: {value: number('SE value', 0), label: "Failed"},
        sw: {value: number('SW value', 0), label: "Succeeded"}
      }
    }
  },
  template: `
      <div class="main-container">
          <div class="content-container" style="width: 300px; height: 200px;">
              <div class="content-area">
                <app-view-quadrant [view]="view"></app-view-quadrant>
              </div>
           </div>
       </div>
 `,
}));
