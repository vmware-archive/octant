import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { LabelsView } from 'src/app/models/content';
import { LabelFilterService } from 'src/app/services/label-filter/label-filter.service';
import { ViewUtil } from 'src/app/util/view';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

@Component({
  selector: 'app-view-labels',
  templateUrl: './labels.component.html',
  styleUrls: ['./labels.component.scss'],
})
export class LabelsComponent implements OnChanges {
  @Input() view: LabelsView;
  title: string;
  labelKeys: string[];
  labels: { [key: string]: string };
  trackByIdentity = trackByIdentity;

  constructor(private labelFilter: LabelFilterService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as LabelsView;

      const vu = new ViewUtil(view);
      this.title = vu.titleAsText();

      this.labels = view.config.labels;
      this.labelKeys = Object.keys(this.labels);
    }
  }

  click(key: string, value: string) {
    this.labelFilter.add({ key, value });
  }
}
