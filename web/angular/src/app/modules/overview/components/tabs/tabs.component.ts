import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { View } from 'src/app/models/content';
import { ViewUtil } from 'src/app/util/view';

interface Tab {
  name: string;
  active: boolean;
  view: View;
  accessor: string;
}

@Component({
  selector: 'app-object-tabs',
  templateUrl: './tabs.component.html',
  styleUrls: ['./tabs.component.scss'],
})
export class TabsComponent implements OnChanges {
  @Input() title: string;
  @Input() views: View[];

  tabs: Tab[] = [];
  checkTab = false;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.views.currentValue) {
      const views = changes.views.currentValue as View[];
      this.tabs = views.map((view, index) => {
        const vu = new ViewUtil(view);
        const title = vu.titleAsText();

        return {
          name: title,
          active: index === 0,
          view,
          accessor: view.metadata.accessor,
        };
      });

      this.checkTab = true;
    }
  }
}
