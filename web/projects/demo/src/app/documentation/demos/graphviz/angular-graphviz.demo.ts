import { Component } from '@angular/core';
import { GraphvizView } from '../../../../../../../src/app/modules/shared/models/content';

const view: GraphvizView = {
  config: {
    dot: 'test',
  },
  metadata: {
    type: 'graphviz',
  },
};

const code = `graphviz component
`;

@Component({
  selector: 'app-angular-graphviz-demo',
  templateUrl: './angular-graphviz.demo.html',
})
export class AngularGraphvizDemoComponent {
  view = view;
  code = code;
}
