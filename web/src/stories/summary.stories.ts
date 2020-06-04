import { storiesOf } from '@storybook/angular';
import {SummaryComponent} from "../app/modules/shared/components/presentation/summary/summary.component";
import {object} from "@storybook/addon-knobs";

storiesOf('Components', module).add('Summary', () => ({
  props: {
    view: object('View',{
      metadata: {
        type: "summary",
        title: [{metadata: {type: "text"}, "config": {value: "Configuration"}}]
      },
      config: {
        sections: [{
          header: "Type",
          content: {metadata: {type: "text"}, "config": {value: "kubernetes.io/service-account-token"}}
        }]
      }
    }),
  },
  component: SummaryComponent,
}));
