import { Component } from '@angular/core';
import { IFrameView } from '../../../../../../../src/app/modules/shared/models/content';

const view: IFrameView = {
  config: {
    url: 'https://octant.dev',
    title: 'title',
  },
  metadata: {
    type: 'iframe',
  },
};

const code = `component.NewIFrame("https://github.com/vmware-tanzu/octant", "title")
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-iframe-demo',
  templateUrl: './angular-iframe.demo.html',
})
export class AngularIFrameDemoComponent {
  view = view;
  code = code;
  json = json;
}
