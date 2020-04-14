import { Component } from '@angular/core';
import { TerminalView } from '../../../../../../../src/app/modules/shared/models/content';

const view: TerminalView = {
  config: {
    name: 'example',
    namespace: 'default',
    podName: 'nginx-pod',
    containers: ['nginx', 'nginx-sidecar'],
    terminal: {
      container: 'nginx',
      command: '/bin/sh',
      active: false,
    },
  },
  metadata: {
    type: 'terminal',
  },
};

const code = `terminal component
`;

@Component({
  selector: 'app-angular-terminal-demo',
  templateUrl: './angular-terminal.demo.html',
})
export class AngularTerminalDemoComponent {
  view = view;
  code = code;
}
