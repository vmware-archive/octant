import { Component } from '@angular/core';
import { LinkView } from '../../../../../../../src/app/modules/shared/models/content';

const view: LinkView = {
  config: {
    ref: 'https://example.org/',
    value: 'Example Link',
  },
  metadata: {
    type: 'link',
  },
};

const code = `title = "Example Link"
ref = "https://example.org"
value = "Example Link"
component.NewLink(title, value, ref)
`;
const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-link-demo',
  templateUrl: './angular-link.demo.html',
})
export class AngularLinkDemoComponent {
  view = view;
  code = code;
  json = json;
}
