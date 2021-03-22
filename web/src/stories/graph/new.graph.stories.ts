import { storiesOf } from '@storybook/angular';
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
import { object } from '@storybook/addon-knobs';

const testCases = [
  { title: 'Deployment', data: REAL_DATA_DEPLOYMENT },
  { title: 'StatefulSet', data: REAL_DATA_STATEFUL_SET },
  { title: 'DaemonSet', data: REAL_DATA_DAEMON_SET },
  { title: 'single DaemonSet', data: REAL_DATA_DAEMON_SET2 },
  { title: 'two ReplicaSets', data: REAL_DATA_TWO_REPLICAS },
  { title: 'Job', data: REAL_DATA_JOB },
  { title: 'Ingress', data: REAL_DATA_INGRESS },
  { title: 'CRDs', data: REAL_DATA_CRDS },
  { title: 'more CRDs', data: REAL_DATA_CRDS2 },
];

testCases.map(story =>
  storiesOf('Other/Resources', module).add(`with ${story.title}`, () => {
    const eles = object('elements', { config: story.data });

    return {
      props: {
        elements: eles,
      },
      template: `
        <div class="main-container">
            <div class="content-container">
                <div class="content-area" style="background-color: white; padding: 0px 16px 0px 0px;">
                    <app-view-resource-viewer [view]="elements">
                    </app-view-resource-viewer>
                </div>
            </div>
        </div>
        `,
    };
  })
);
