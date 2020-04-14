import { Component } from '@angular/core';
import { PortForwardView } from '../../../../../../../src/app/modules/shared/models/content';

const view: PortForwardView = {
  config: {
    text: '',
    id: '',
    action: '',
    status: '',
    ports: [],
    target: null,
  },
  metadata: {
    type: 'selectors',
  },
};

const code = `selector component
`;

@Component({
  selector: 'app-angular-portforward-demo',
  templateUrl: './angular-portforward.demo.html',
})
export class AngularPortForwardDemoComponent {
  view = view;
  code = code;
}
