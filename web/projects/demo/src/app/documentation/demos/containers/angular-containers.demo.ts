import { Component } from '@angular/core';
import {
  ContainersView,
  ContainerDef,
} from '../../../../../../../src/app/modules/shared/models/content';

const container: ContainerDef = {
  name: 'nginx-container',
  image: 'nginx:1.17.0',
};

const view: ContainersView = {
  config: {
    containers: [container],
  },
  metadata: {
    type: 'containers',
  },
};

const code = `containers code`;

@Component({
  selector: 'app-angular-containers-demo',
  templateUrl: './angular-containers.demo.html',
})
export class AngularContainersDemoComponent {
  view = view;
  code = code;
}
