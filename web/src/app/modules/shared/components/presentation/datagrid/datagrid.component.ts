// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import {
  TableFilters,
  TableRow,
  TableView,
  View,
} from 'src/app/modules/shared/models/content';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { ViewService } from '../../../services/view/view.service';

@Component({
  selector: 'app-view-datagrid',
  templateUrl: './datagrid.component.html',
  styleUrls: ['./datagrid.component.scss'],
})
export class DatagridComponent implements OnChanges {
  private v: TableView;

  @Input() set view(v: View) {
    this.v = v as TableView;
  }
  get view() {
    return this.v;
  }

  columns: string[];
  rows: TableRow[];
  title: string;
  placeholder: string;
  lastUpdated: Date;
  filters: TableFilters;

  private previousView: SimpleChanges;

  identifyRow = trackByIndex;
  identifyColumn = trackByIdentity;
  loading: boolean;

  constructor(private viewService: ViewService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      if (
        JSON.stringify(this.previousView) !==
        JSON.stringify(changes.view.currentValue)
      ) {
        this.title = this.viewService.viewTitleAsText(this.view);

        const current = changes.view.currentValue;
        this.columns = current.config.columns.map(column => column.name);
        this.rows = current.config.rows;
        this.placeholder = current.config.emptyContent;
        this.lastUpdated = new Date();
        this.loading = current.config.loading;
        this.filters = current.config.filters;

        this.previousView = changes.view.currentValue;
      }
    }
  }
}
