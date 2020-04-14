import { Component } from '@angular/core';
import { PodStatusView } from '../../../../../../../src/app/modules/shared/models/content';

const view: PodStatusView = {
  config: {
    pods: {},
  },
  metadata: {
    type: 'podStatus',
  },
};

const code = `pod status component
`;

@Component({
  selector: 'app-angular-pod-status-demo',
  templateUrl: './angular-pod-status.demo.html',
})
export class AngularPodStatusDemoComponent {
  view = view;
  code = code;
}
