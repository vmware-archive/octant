import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { FlexLayoutItem, FlexLayoutView } from 'src/app/models/content';

@Component({
  selector: 'app-view-flexlayout',
  templateUrl: './flexlayout.component.html',
  styleUrls: ['./flexlayout.component.scss'],
})
export class FlexlayoutComponent implements OnChanges {
  @Input()
  view: FlexLayoutView;

  sections: FlexLayoutItem[][];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as FlexLayoutView;
      this.sections = view.config.sections;
    }
  }
}
