import { big_data } from './overview.data';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { argTypesView } from './helpers/helpers';

export default {
  title: 'Other/Sandbox',
} as Meta;

export const sandboxStory: Story = args => {
  return {
    props: {
      view: args.view,
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
};

sandboxStory.storyName = 'Component Sandbox';

sandboxStory.argTypes = argTypesView;

sandboxStory.args = {
  view: big_data,
};
