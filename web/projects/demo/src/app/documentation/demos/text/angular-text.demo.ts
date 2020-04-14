import { Component } from '@angular/core';
import { TextView } from '../../../../../../../src/app/modules/shared/models/content';

const sampleText = 'This is a text';

const view: TextView = {
  config: {
    value: sampleText,
  },
  metadata: {
    type: 'text',
  },
};

const code = `text := component.NewText("${sampleText}")
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-text-demo',
  templateUrl: './angular-text.demo.html',
})
export class AngularTextDemoComponent {
  view = view;
  code = code;
  json = json;
}
