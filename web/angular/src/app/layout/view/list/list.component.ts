import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ListView, View } from 'src/app/models/content';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ViewListComponent implements OnChanges {
  @Input() listView: ListView;

  items: View[];

  constructor() {}

  rawList: string;

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.listView.currentValue) {
      const view = changes.listView.currentValue as ListView;
      this.items = view.config.items;
    }
  }
}
