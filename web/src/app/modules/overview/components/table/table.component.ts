import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TableRow, TableView } from 'src/app/models/content';
import { ViewUtil } from 'src/app/util/view';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';

@Component({
  selector: 'app-view-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
})
export class TableComponent implements OnChanges {
  @Input() view: TableView;
  columns: string[];
  rows: TableRow[];
  title: string;
  placeholder: string;
  trackByIdentity = trackByIdentity;
  trackByIndex = trackByIndex;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const current = changes.view.currentValue;
      const vu = new ViewUtil(current);
      this.title = vu.titleAsText();
      this.columns = current.config.columns.map((column) => column.name);
      this.rows = current.config.rows;
      this.placeholder = current.config.emptyContent;
    }
  }
}
