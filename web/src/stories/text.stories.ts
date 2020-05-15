import { storiesOf } from '@storybook/angular';
import { boolean, object } from '@storybook/addon-knobs';

storiesOf('Text', module).add('with text', () => ({
  props: {
    view: {
      metadata: {
        type: 'text',
      },
      config: {
        value: object('value', 'hello world'),
        isMarkdown: boolean('isMarkdown', false),
      },
    },
  },
  template: `
    <div class="main-container">
        <div class="content-container">
            <div class="content-area">
                <app-view-text [view]="view"></app-view-text>
            </div>
        </div>
    </div>
    `,
}));
