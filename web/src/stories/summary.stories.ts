import { storiesOf } from '@storybook/angular';
import {SummaryComponent} from "../app/modules/shared/components/presentation/summary/summary.component";
import {text} from "@storybook/addon-knobs";

storiesOf('Summary', module).add('Summary component', () => ({
  props: {
    view: {
      metadata: {
        type: "summary",
        title: [{metadata: {type: "text"}, "config": {value: text("Title", "Configuration")}}]
      },
      config: {
        sections: [{
          header: text("Header", "Type"),
          content: {metadata: {type: "text"}, "config": {value: text("Value","kubernetes.io/service-account-token")}}
        }]
      }
    },
  },
  component: SummaryComponent,
}));
