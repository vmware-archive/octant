// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TableRow, TableView } from 'src/app/shared/models/content';
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
  placeholder: string;
  trackByIdentity = trackByIdentity;
  trackByIndex = trackByIndex;

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const current = changes.view.currentValue;
      this.columns = current.config.columns.map(column => column.name);
      this.rows = current.config.rows;
      this.placeholder = current.config.emptyContent;
    }
  }
}
