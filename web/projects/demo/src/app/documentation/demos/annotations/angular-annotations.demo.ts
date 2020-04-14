import { Component } from '@angular/core';
import {
  View,
  AnnotationsView,
  TextView,
  SummaryView,
} from '../../../../../../../src/app/modules/shared/models/content';

const view: AnnotationsView = {
  config: {
    annotations: {
      ['foo']: 'bar',
    },
  },
  metadata: {
    type: 'annotations',
  },
};

const code1 = `annotations := map[string]string{"foo": "bar"}
component.NewAnnotations(annotations)`;

const text: TextView = {
  config: {
    value: 'Card Title',
  },
  metadata: {
    type: 'text',
  },
};

const summaryView: SummaryView = {
  metadata: {
    type: 'summary',
    title: [text],
  },
  config: {
    sections: [
      {
        header: 'Annotations',
        content: view,
      },
    ],
    actions: [],
  },
};

const code2 = `component.NewSummary("Card Title", []component.SummarySections{
    {
        Header: "Annotations",
        Content: component.NewAnnotations(map[string]string{"foo": "bar"}),
    },
})`;

const json1 = JSON.stringify(view, null, 4);
const json2 = JSON.stringify(summaryView, null, 4);
@Component({
  selector: 'app-angular-annotations-demo',
  templateUrl: './angular-annotations.demo.html',
})
export class AngularAnnotationsDemoComponent {
  view = view;
  summaryView = summaryView;
  preview = view as View;
  code1 = code1;
  json1 = json1;
  code2 = code2;
  json2 = json2;
}
