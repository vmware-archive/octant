import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View } from 'src/app/models/content';

@Component({
  selector: 'app-content-switcher',
  templateUrl: './content-switcher.component.html',
  styleUrls: ['./content-switcher.component.scss'],
})
export class ContentSwitcherComponent implements OnChanges {
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
