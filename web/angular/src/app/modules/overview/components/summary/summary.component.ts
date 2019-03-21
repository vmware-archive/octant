import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { SummaryItem, SummaryView } from 'src/app/models/content';
import { ViewUtil } from 'src/app/util/view';

@Component({
  selector: 'app-view-summary',
  templateUrl: './summary.component.html',
  styleUrls: ['./summary.component.scss'],
})
export class SummaryComponent implements OnChanges {
  @Input() view: SummaryView;

  title: string;

  items: SummaryItem[];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue;

      const vu = new ViewUtil(view);
      this.title = vu.titleAsText();
      this.items = view.config.sections;
    }
  }
}
