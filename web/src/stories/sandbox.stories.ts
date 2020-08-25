import { storiesOf } from '@storybook/angular';
import { object } from '@storybook/addon-knobs';
import { big_data } from './overview.data';

storiesOf('Other/Sandbox', module).add('Component Sandbox', () => {
  const view = object('JSON', big_data);

  return {
    props: {
      view,
    },
    template: `
      <div class="main-container">
          <div class="content-container">
              <div class="content-area">
                <app-view-container [view]="view">
                </app-view-container>
              </div>
          </div>
      </div>
      `,
  };
});
