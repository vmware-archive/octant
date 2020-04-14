import { Component } from '@angular/core';
import { LogsView } from '../../../../../../../src/app/modules/shared/models/content';

const view: LogsView = {
  config: {
    namespace: 'default',
    name: 'nginx-deployment-75c5cb5f44-7wh55',
    containers: ['[all-containers]', 'nginx'],
  },
  metadata: {
    type: 'logs',
  },
};

const code = `component.NewLogs(pod.Namespace, pod.Name, containerNames...)
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-logs-demo',
  templateUrl: './angular-logs.demo.html',
})
export class AngularLogsDemoComponent {
  view = view;
  code = code;
  json = json;
}
