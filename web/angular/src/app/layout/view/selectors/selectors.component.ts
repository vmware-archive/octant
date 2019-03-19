import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ExpressionSelectorView, LabelSelectorView, SelectorsView } from 'src/app/models/content';

@Component({
  selector: 'app-view-selectors',
  templateUrl: './selectors.component.html',
  styleUrls: ['./selectors.component.scss'],
})
export class SelectorsComponent implements OnChanges {
  @Input() view: SelectorsView;

  selectors: Array<ExpressionSelectorView | LabelSelectorView> = [];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as SelectorsView;
      this.selectors = view.config.selectors;
    }
  }
}
