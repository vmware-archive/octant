import { Component } from '@angular/core';
import { ExpressionSelectorView } from '../../../../../../../src/app/modules/shared/models/content';

const view: ExpressionSelectorView = {
  config: {
    key: 'foo',
    operator: 'bar',
    values: ['a', 'b', 'c'],
  },
  metadata: {
    type: 'expressionSelector',
  },
};

const code = `expression selector component
`;

@Component({
  selector: 'app-angular-expression-selector-demo',
  templateUrl: './angular-expression-selector.demo.html',
})
export class AngularExpressionSelectorDemoComponent {
  view = view;
  code = code;
}
