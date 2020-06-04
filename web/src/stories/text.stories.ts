import { storiesOf } from '@storybook/angular';
import { boolean, object } from '@storybook/addon-knobs';
import {TextComponent} from "../app/modules/shared/components/presentation/text/text.component";

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
  component: TextComponent,
}));
