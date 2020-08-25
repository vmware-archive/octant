// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import { TableRow, TableView } from 'src/app/modules/shared/models/content';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-table',
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
})
export class TableComponent extends AbstractViewComponent<TableView> {
  columns: string[];
  rows: TableRow[];
  title: string;
  placeholder: string;
  trackByIdentity = trackByIdentity;
  trackByIndex = trackByIndex;

  constructor(private viewService: ViewService) {
    super();
  }

  update() {
    const current = this.v;
    this.title = this.viewService.viewTitleAsText(current);
    this.columns = current.config.columns.map(column => column.name);
    this.rows = current.config.rows;
    this.placeholder = current.config.emptyContent;
  }
}
