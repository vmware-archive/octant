import { Component } from '@angular/core';
import {
  QuadrantView,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const title: TextView = {
  config: {
    value: 'Summary',
  },
  metadata: {
    type: 'text',
  },
};

const view: QuadrantView = {
  config: {
    nw: {
      value: '3',
      label: 'Running',
    },
    ne: {
      value: '0',
      label: 'Waiting',
    },
    se: {
      value: '0',
      label: 'Failed',
    },
    sw: {
      value: '0',
      label: 'Succeeded',
    },
  },
  metadata: {
    type: 'quadrant',
    title: [title],
  },
};

const code = `quadrant := component.NewQuadrant("Status")
quadrant.Set(component.QuadNW, "Running", fmt.Sprintf("%d", ps.Running))
quadrant.Set(component.QuadNE, "Waiting", fmt.Sprintf("%d", ps.Waiting))
quadrant.Set(component.QuadSE, "Failed", fmt.Sprintf("%d", ps.Failed))
quadrant.Set(component.QuadSW, "Succeeded", fmt.Sprintf("%d", ps.Succeeded))
`;
const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-quadrant-demo',
  templateUrl: './angular-quadrant.demo.html',
})
export class AngularQuadrantDemoComponent {
  view = view;
  code = code;
  json = json;
}
