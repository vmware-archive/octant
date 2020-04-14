import { Component } from '@angular/core';
import {
  FlexLayoutView,
  FlexLayoutItem,
  SummaryView,
  SummaryItem,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const text: TextView = {
  config: {
    value: 'sample text',
  },
  metadata: {
    type: 'text',
  },
};

const summaryItemFull: SummaryItem = {
  header: 'Full Width',
  content: text,
};

const summaryItemHalf: SummaryItem = {
  header: 'Half Width',
  content: text,
};

const summaryItemQuarter: SummaryItem = {
  header: 'Quarter Width',
  content: text,
};

const summaryFull: SummaryView = {
  config: {
    sections: [summaryItemFull],
    actions: [],
  },
  metadata: {
    type: 'summary',
  },
};

const summaryHalf: SummaryView = {
  config: {
    sections: [summaryItemHalf],
    actions: [],
  },
  metadata: {
    type: 'summary',
  },
};

const summaryQuarter: SummaryView = {
  config: {
    sections: [summaryItemQuarter],
    actions: [],
  },
  metadata: {
    type: 'summary',
  },
};

const sectionFull: FlexLayoutItem = {
  width: 24,
  height: null,
  view: summaryFull,
};

const sectionHalf: FlexLayoutItem = {
  width: 12,
  height: null,
  view: summaryHalf,
};

const sectionQuarter: FlexLayoutItem = {
  width: 6,
  height: null,
  view: summaryQuarter,
};

const view: FlexLayoutView = {
  config: {
    sections: [[sectionFull, sectionHalf, sectionQuarter]],
    buttonGroup: null,
  },
  metadata: {
    type: 'flexlayout',
  },
};

const code1 = `flexLayout := component.NewFlexLayout("Summary")
flexLayout.AddSections([]component.FlexLayoutSection{
  {
    {
      Width: component.WidthFull,
      View: component.NewSummary("Metadata", component.SummarySections{
        {
          Header:  "Full Width",
          Content: component.NewText("sample text"),
        },
      }...),
    },
    {
      Width: component.WidthHalf,
      View: component.NewSummary("Metadata", component.SummarySections{
        {
          Header:  "Half Width",
          Content: component.NewText("sample text"),
        },
      }...),
    },
    {
      Width: component.WidthQuarter,
      View: component.NewSummary("Metadata", component.SummarySections{
        {
          Header:  "Quarter Width",
          Content: component.NewText("sample text"),
        },
      }...),
    },
  },
}...)`;

const json1 = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-flexlayout-demo',
  templateUrl: './angular-flexlayout.demo.html',
})
export class AngularFlexLayoutDemoComponent {
  view = view;
  code1 = code1;
  json1 = json1;
}
