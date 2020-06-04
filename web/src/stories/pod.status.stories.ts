import {PodStatusView, TextView} from "../app/modules/shared/models/content";
import {storiesOf} from "@storybook/angular";
import {PodStatusComponent} from "../app/modules/shared/components/presentation/pod-status/pod-status.component";
import {object} from "@storybook/addon-knobs";

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
      "pod-57466fd965-xprw9": {
        "details": [statusText], "status": "ok"
      }
    }
  }
};

storiesOf('Components', module).add('Pod Status', () => ({
  props: {
    view: object('View', view)
  },
  component: PodStatusComponent,
}));
