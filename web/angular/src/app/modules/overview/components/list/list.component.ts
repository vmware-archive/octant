import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { ListView, View } from 'src/app/models/content';
import isEqual from 'lodash.isequal';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ViewListComponent implements OnChanges {
  @Input() listView: ListView;
  items: View[];

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.listView.currentValue) {
      const view = changes.listView.currentValue as ListView;
      // note(marlon): angular runs a dom update every time
      // a new object is set to `this.items`. this means that
      // every time we poll on the event stream and return a
      // new object ref the dom refreshes regardless if the
      // actual data within the object ref has changed. this
      // should change with angular ivy (upcoming release)
      if (!isEqual(this.items, view.config.items)) {
        this.items = view.config.items;
      }
    }
  }
}
