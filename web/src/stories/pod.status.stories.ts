import {PodStatusView, TextView} from "../app/modules/shared/models/content";
import {storiesOf} from "@storybook/angular";
import {PodStatusComponent} from "../app/modules/shared/components/presentation/pod-status/pod-status.component";

const statusText: TextView = {
  config: {
    value: '',
  },
  metadata: {
    type: 'text',
  },
};

const view: PodStatusView = {
  "metadata": {"type": "podStatus"},
  "config": {
    "pods": {
      "coreapi-57466fd965-xprw9": {
        "details": [statusText], "status": "ok"
      }
    }
  }
};

storiesOf('Pod Status', module).add('Pod Status component', () => ({
  props: {
    view: view
  },
  component: PodStatusComponent,
}));
