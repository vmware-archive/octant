import { storiesOf } from '@storybook/angular';
import {TimestampComponent} from "../app/modules/shared/components/presentation/timestamp/timestamp.component";
import {object} from "@storybook/addon-knobs";

storiesOf('Components', module).add('Timestamp', () => ({
  props: {
    view: object('View', {metadata: {type: "timestamp"}, config: {timestamp: 1588716648}}),
  },
  component: TimestampComponent,
}));
