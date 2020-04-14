import { Component } from '@angular/core';
import {
  LinkView,
  SummaryView,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const title: TextView = {
  config: {
    value: 'Configuration',
  },
  metadata: {
    type: 'text',
  },
};

const content1: TextView = {
  config: {
    value: '0',
  },
  metadata: {
    type: 'text',
  },
};

const content2: LinkView = {
  config: {
    value: 'minikube',
    ref: '/cluster-overview/nodes/minikube',
  },
  metadata: {
    type: 'link',
  },
};

const content3: LinkView = {
  config: {
    value: 'default',
    ref:
      '/overview/namespace/default/config-and-storage/service-accounts/default',
  },
  metadata: {
    type: 'link',
  },
};

const view: SummaryView = {
  config: {
    sections: [
      {
        header: 'Priority',
        content: content1,
      },
      {
        header: 'Node',
        content: content2,
      },
      {
        header: 'Service Account',
        content: content3,
      },
    ],
    actions: [],
  },
  metadata: {
    type: 'summary',
    title: [title],
  },
};

const code = `sections := component.SummarySections{}
sections.AddText("Priority", fmt.Sprintf("%d", *pod.Spec.Priority))
sections = append(sections, component.SummarySection{
  Header:  "Service Account",
  Content: contentLink,
})
sections.Add("Node", nodeLink)
summary := component.NewSummary("Configuration", sections...)
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-summary-demo',
  templateUrl: './angular-summary.demo.html',
})
export class AngularSummaryDemoComponent {
  view = view;
  code = code;
  json = json;
}
