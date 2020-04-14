import { Component } from '@angular/core';
import { SelectorsView } from '../../../../../../../src/app/modules/shared/models/content';

const view: SelectorsView = {
  config: {
    selectors: [
      {
        metadata: {
          type: 'labelSelector',
        },
        config: {
          key: 'app',
          value: 'httpbin',
        },
      },
    ],
  },
  metadata: {
    type: 'selectors',
  },
};

const code = `selectors:= component.NewSelectors([]component.Selector{component.NewLabelSelector("app", "theapp")})
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-selectors-demo',
  templateUrl: './angular-selectors.demo.html',
})
export class AngularSelectorsDemoComponent {
  view = view;
  code = code;
  json = json;
}
