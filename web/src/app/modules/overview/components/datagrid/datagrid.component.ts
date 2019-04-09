import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TableRow, TableView } from 'src/app/models/content';
import { ViewUtil } from 'src/app/util/view';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';

@Component({
  selector: 'app-view-datagrid',
  templateUrl: './datagrid.component.html',
  styleUrls: ['./datagrid.component.scss'],
})
export class DatagridComponent implements OnChanges {
  @Input() view: TableView;

  columns: string[];
  rows: TableRow[];
  title: string;
  placeholder: string;

  lastUpdated: Date;

  identifyRow = trackByIndex;
  identifyColumn = trackByIdentity;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const vu = new ViewUtil(this.view);
      this.title = vu.titleAsText();

      const current = changes.view.currentValue;
      this.columns = current.config.columns.map((column) => column.name);
      this.rows = current.config.rows;
      this.placeholder = current.config.emptyContent;
      this.lastUpdated = new Date();
    }
  }

}
