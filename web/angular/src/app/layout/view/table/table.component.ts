import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TableRow, TableView } from 'src/app/models/content';
import { ViewUtil } from 'src/app/util/view';

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

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.tableView) {
      const vu = new ViewUtil(this.view);
      this.title = vu.titleAsText();

      const current = changes.tableView.currentValue;
      this.columns = current.config.columns.map((column) => column.name);
      this.rows = current.config.rows;
      this.placeholder = current.config.emptyContent;
    }
  }
}
