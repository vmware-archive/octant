import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View } from 'src/app/models/content';

@Component({
  selector: 'app-view',
  templateUrl: './view.component.html',
  styleUrls: ['./view.component.scss'],
})
export class ViewComponent implements OnChanges {
  @Input() view: View;

  currentViewType: string;
  currentView: View;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as View;

      if (view.metadata) {
        this.currentViewType = view.metadata.type;
        this.currentView = view;

        this.view = view;
      }
    }
  }
}
