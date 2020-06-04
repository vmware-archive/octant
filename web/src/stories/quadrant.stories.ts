import { storiesOf } from '@storybook/angular';
import {object} from '@storybook/addon-knobs';


storiesOf('Components', module).add('Quadrant', () => ({
  props: {
    view: object('View',{
      metadata: {
        type: "quadrant",
        title: [{metadata: {type: "text"}, config: {value: "Status"}}]
      },
      config: {
        nw: {value: 1, label: "Running"},
        ne: {value: 0, label: "Waiting"},
        se: {value: 0, label: "Failed"},
        sw: {value: 0, label: "Succeeded"}
      }
    })
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
