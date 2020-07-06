import {storiesOf} from '@storybook/angular';
import {object} from "@storybook/addon-knobs";
import {big_data} from "./overview.data";

storiesOf('Sandbox', module).add('Component Sandbox', () => {
  const view = object('JSON', big_data);

  return {
    props: {
      view: view,
    },
    template: `
      <div class="main-container">
          <div class="content-container">
              <div class="content-area">
                <app-content-switcher [view]="view">
                </app-content-switcher>
              </div>
          </div>
      </div>
      `,
  }
});
