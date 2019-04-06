import { Component, Input } from '@angular/core';
import { ListView, View } from 'src/app/models/content';
import { titleAsText } from 'src/app/util/view';

@Component({
  selector: 'app-view-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss'],
})
export class ViewListComponent {
  @Input() listView: ListView;

  identifyItem(index: number, item: View): string {
    return titleAsText(item.metadata.title);
  }
}
