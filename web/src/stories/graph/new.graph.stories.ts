import {
  REAL_DATA_DAEMON_SET,
  REAL_DATA_DAEMON_SET2,
  REAL_DATA_DEPLOYMENT,
  REAL_DATA_INGRESS,
  REAL_DATA_JOB,
  REAL_DATA_STATEFUL_SET,
  REAL_DATA_TWO_REPLICAS,
  REAL_DATA_CRDS,
  REAL_DATA_CRDS2,
} from './graph.real.data';
import { Meta, Story } from '@storybook/angular/types-6-0';
import { ResourceViewerComponent } from '../../app/modules/shared/components/presentation/resource-viewer/resource-viewer.component';
import { argTypesView } from '../helpers/helpers';

export default {
  title: 'Other/Resources',
} as Meta;

const Template: Story<ResourceViewerComponent> = args => ({
  props: {
    view: { config: args.view },
  },
  argTypes: argTypesView,
  template: `
        <div class="main-container">
            <div class="content-container">
                <div class="content-area" style="background-color: white; padding: 0px 16px 0px 0px;">
                    <app-view-resource-viewer [view]="view">
                    </app-view-resource-viewer>
                </div>
            </div>
        </div>
        `,
});

export const graphStoryDeployment = Template.bind({});
graphStoryDeployment.storyName = 'Deployment';
graphStoryDeployment.args = {
  view: REAL_DATA_DEPLOYMENT,
};

export const graphStoryStatefulSet = Template.bind({});
graphStoryStatefulSet.storyName = 'Stateful Set';
graphStoryStatefulSet.args = {
  view: REAL_DATA_STATEFUL_SET,
};

export const graphStoryDaemonSet = Template.bind({});
graphStoryDaemonSet.storyName = 'Daemon Set';
graphStoryDaemonSet.args = {
  view: REAL_DATA_DAEMON_SET,
};

export const graphStoryDaemonSet2 = Template.bind({});
graphStoryDaemonSet2.storyName = 'Single DaemonSet';
graphStoryDaemonSet2.args = {
  view: REAL_DATA_DAEMON_SET2,
};

export const graphStoryReplicas = Template.bind({});
graphStoryReplicas.storyName = 'Two ReplicaSets';
graphStoryReplicas.args = {
  view: REAL_DATA_TWO_REPLICAS,
};

export const graphStoryJob = Template.bind({});
graphStoryJob.storyName = 'Job';
graphStoryJob.args = {
  view: REAL_DATA_JOB,
};

export const graphStoryIngress = Template.bind({});
graphStoryIngress.storyName = 'Ingress';
graphStoryIngress.args = {
  view: REAL_DATA_INGRESS,
};

export const graphStoryCRDs = Template.bind({});
graphStoryCRDs.storyName = 'CRDs';
graphStoryCRDs.args = {
  view: REAL_DATA_CRDS,
};

export const graphStoryCRDs2 = Template.bind({});
graphStoryCRDs2.storyName = 'More CRDs';
graphStoryCRDs2.args = {
  view: REAL_DATA_CRDS2,
};
