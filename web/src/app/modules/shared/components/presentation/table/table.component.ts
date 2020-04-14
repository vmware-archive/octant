// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { TableRow, TableView } from '../../../../shared/models/content';
import trackByIdentity from '../../../../../util/trackBy/trackByIdentity';
import trackByIndex from '../../../../../util/trackBy/trackByIndex';
import { ViewService } from '../../../services/view/view.service';

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

  constructor(private viewService: ViewService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      const current = changes.view.currentValue;
      this.title = this.viewService.viewTitleAsText(current);
      this.columns = current.config.columns.map(column => column.name);
      this.rows = current.config.rows;
      this.placeholder = current.config.emptyContent;
    }
  }
}
