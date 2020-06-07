import { storiesOf } from '@storybook/angular';
import {TimestampComponent} from "../app/modules/shared/components/presentation/timestamp/timestamp.component";

storiesOf('Timestamp', module).add('Timestamp component', () => ({
  props: {
    view: {metadata: {type: "timestamp"}, config: {timestamp: 1588716648}},
  },
  component: TimestampComponent,
}));
