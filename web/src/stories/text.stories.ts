import { storiesOf } from '@storybook/angular';
import { object } from '@storybook/addon-knobs';
import {TextComponent} from "../app/modules/shared/components/presentation/text/text.component";

storiesOf('Components', module).add('Text', () => ({
  props: {
    view: object('View',{
      metadata: {
        type: 'text',
      },
      config: {
        value: 'hello world',
        isMarkdown: false,
      },
    }),
  },
  component: TextComponent,
}));
