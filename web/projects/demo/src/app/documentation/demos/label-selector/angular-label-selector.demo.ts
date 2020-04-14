import { Component } from '@angular/core';
import { LabelSelectorView } from '../../../../../../../src/app/modules/shared/models/content';

const view: LabelSelectorView = {
  config: {
    key: 'foo',
    value: 'bar',
  },
  metadata: {
    type: 'labelSelector',
  },
};

const code = `label selector component
`;

@Component({
  selector: 'app-angular-label-selector-demo',
  templateUrl: './angular-label-selector.demo.html',
})
export class AngularLabelSelectorDemoComponent {
  view = view;
  code = code;
}
