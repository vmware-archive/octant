import { Component } from '@angular/core';
import {
  LabelsView,
  LabelSelectorView,
  SelectorsView,
  ExpressionSelectorView,
} from '../../../../../../../src/app/modules/shared/models/content';

const labelView: LabelsView = {
  config: {
    labels: { ['foo']: 'bar' },
  },
  metadata: {
    type: 'labels',
  },
};

const selectorsView: SelectorsView = {
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

const expressionSelectorView: ExpressionSelectorView = {
  config: {
    key: 'key',
    operator: 'In',
    values: ['a', 'b'],
  },
  metadata: {
    type: 'expressionSelector',
  },
};

const labelCode = `labels = map[string]string{"foo": "bar"}
component.NewLabels(labels)`;
const labelJSON = JSON.stringify(labelView, null, 4);

const selectorsCode = `component.NewSelectors([]component.Selector{
    component.NewLabelSelector("app", "httpbin")
})
`;
const selectorsJSON = JSON.stringify(selectorsView, null, 4);

const expressionSelectorCode = `component.NewExpressionSelector("key", component.OperatorIn, []string{"a", "b"})`;
const expressionSelectorJSON = JSON.stringify(expressionSelectorView, null, 4);

@Component({
  selector: 'app-angular-labels-demo',
  templateUrl: './angular-labels.demo.html',
})
export class AngularLabelsDemoComponent {
  labelView = labelView;
  labelCode = labelCode;
  labelJSON = labelJSON;
  selectorsView = selectorsView;
  selectorsCode = selectorsCode;
  selectorsJSON = selectorsJSON;
  expressionSelectorView = expressionSelectorView;
  expressionSelectorCode = expressionSelectorCode;
  expressionSelectorJSON = expressionSelectorJSON;
}
