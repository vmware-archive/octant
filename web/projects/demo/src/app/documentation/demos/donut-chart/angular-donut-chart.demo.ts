import { Component } from '@angular/core';
import {
  DonutChartView,
  FlexLayoutView,
  FlexLayoutItem,
  DonutSegment,
  DonutChartLabels,
} from '../../../../../../../src/app/modules/shared/models/content';

const segment: DonutSegment = {
  count: 3,
  status: 'ok',
};

const label: DonutChartLabels = {
  plural: 'Pods',
  singular: 'Pod',
};

const dontChuartView: DonutChartView = {
  config: {
    segments: [segment],
    labels: label,
    size: 3,
  },
  metadata: {
    type: 'donutChart',
  },
};

const sectionQuarter: FlexLayoutItem = {
  width: 6,
  height: null,
  view: dontChuartView,
};

const view: FlexLayoutView = {
  config: {
    sections: [[sectionQuarter]],
    buttonGroup: null,
  },
  metadata: {
    type: 'flexlayout',
  },
};

const code = `component.DonutChart{
	Config: component.DonutChartConfig{
		Size: component.DonutChartSizeMedium,
		Segments: []component.DonutSegment{
			{
				Count: 3,
				Status: component.NodeStatusOK,
			},
		},
		Labels: component.DonutChartLabels{
			Plural: "Pods",
			Singular: "Pod",
		},
	},
}
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-donut-chart-demo',
  templateUrl: './angular-donut-chart.demo.html',
})
export class AngularDonutChartDemoComponent {
  view = view;
  code = code;
  json = json;
}
