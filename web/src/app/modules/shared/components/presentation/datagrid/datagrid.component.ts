// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  ChangeDetectionStrategy,
  Component,
  Input,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import {
  GridAction,
  GridActionsView,
  TableFilters,
  TableRow,
  TableView,
  View,
} from 'src/app/modules/shared/models/content';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import { ActionService } from '../../../services/action/action.service';
import { ViewService } from '../../../services/view/view.service';

@Component({
  selector: 'app-view-datagrid',
  templateUrl: './datagrid.component.html',
  styleUrls: ['./datagrid.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
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
  rowsWithMetadata: TableRowWithMetadata[];
  title: string;
  placeholder: string;
  lastUpdated: Date;
  filters: TableFilters;

  private previousView: SimpleChanges;

  identifyRow = trackByIndex;
  identifyColumn = trackByIdentity;
  loading: boolean;

  constructor(
    private viewService: ViewService,
    private actionService: ActionService
  ) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view) {
      if (
        JSON.stringify(this.previousView) !==
        JSON.stringify(changes.view.currentValue)
      ) {
        this.title = this.viewService.viewTitleAsText(this.view);

        const current = changes.view.currentValue as TableView;
        this.columns = current.config.columns.map(column => column.name);

        if (current.config.rows) {
          this.rowsWithMetadata = this.getRowsWithMetadata(current.config.rows);
        }

        this.placeholder = current.config.emptyContent;
        this.lastUpdated = new Date();
        this.loading = current.config.loading;
        this.filters = current.config.filters;

        this.previousView = changes.view.currentValue;
      }
    }
  }

  private getRowsWithMetadata(rows: TableRow[]) {
    return rows.map(row => {
      let actions: GridAction[] = [];

      if (row.hasOwnProperty('_action')) {
        actions = (row._action as GridActionsView).config.actions;
      }
      return {
        data: row,
        actions,
      };
    });
  }

  runAction(actionPath: string, payload: {}) {
    const update = { ...payload, action: actionPath };
    this.actionService.perform(update);
  }

  showTitle() {
    if (this.view) {
      return this.view.totalItems === undefined || this.view.totalItems > 1;
    }
    return true;
  }
}

interface TableRowWithMetadata {
  data: TableRow;
  actions?: GridAction[];
}
