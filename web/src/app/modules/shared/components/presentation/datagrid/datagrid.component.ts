// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ClrDatagridSortOrder } from '@clr/angular';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
} from '@angular/core';
import {
  Confirmation,
  GridAction,
  GridActionsView,
  TableFilters,
  TableRow,
  TableRowWithMetadata,
  TableView,
} from 'src/app/modules/shared/models/content';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { TimestampComparator } from '../../../../../util/timestamp-comparator';
import { ViewService } from '../../../services/view/view.service';
import { ActionService } from '../../../services/action/action.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { BehaviorSubject, merge, Observable, timer } from 'rxjs';
import { LoadingService } from '../../../services/loading/loading.service';
import { ButtonGroupView } from '../../../models/content';

@Component({
  selector: 'app-view-datagrid',
  templateUrl: './datagrid.component.html',
  styleUrls: ['./datagrid.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DatagridComponent extends AbstractViewComponent<TableView> {
  timeStampComparator = new TimestampComparator();
  sortOrder: ClrDatagridSortOrder = ClrDatagridSortOrder.UNSORTED;

  columns: string[];
  rowsWithMetadata: TableRowWithMetadata[];
  title: string;
  placeholder: string;
  lastUpdated: Date;
  filters: TableFilters;
  buttonGroup?: ButtonGroupView;
  isModalOpen = false;

  actionDialogOptions: ActionDialogOptions = undefined;

  identifyRow = trackByIndex;
  identifyColumn = trackByIdentity;
  identifyAction = trackByIdentity;

  loading: boolean;
  loading$: Observable<boolean>;

  constructor(
    private viewService: ViewService,
    private actionService: ActionService,
    private loadingService: LoadingService,
    private cdr: ChangeDetectorRef
  ) {
    super();
  }

  update() {
    this.title = this.viewService.viewTitleAsText(this.view);

    this.loading = true;

    const done = new BehaviorSubject(false);
    this.loading$ = this.loadingService.withDelay(done, 250, 1000);

    setTimeout(() => {
      this.rowsWithMetadata = this.getRowsWithMetadata(this.v.config.rows);
      this.placeholder = this.v.config.emptyContent;
      this.lastUpdated = new Date();
      this.loading = this.v.config.loading;
      done.next(true);
      done.complete();
      this.cdr.markForCheck();
    });
    this.columns = this.v.config.columns.map(column => column.name);
    this.filters = this.v.config.filters;
    this.buttonGroup = this.v.config.buttonGroup;
  }

  private getRowsWithMetadata(rows: TableRow[]): TableRowWithMetadata[] {
    if (!rows) {
      return [];
    }

    return rows.map(row => {
      let actions: GridAction[] = [];

      if (row.hasOwnProperty('_action')) {
        actions = (row._action as GridActionsView).config.actions;
      }

      const isDeleted = !!row._isDeleted;

      return {
        data: row,
        actions,
        isDeleted,
      };
    });
  }

  runAction(action: GridAction) {
    if (!action.confirmation) {
      const update = { ...action.payload, action: action.actionPath };
      this.actionService.perform(update);
      return;
    }

    this.actionDialogOptions = {
      action,
      text: action.name,
      type: action.type,
      confirmation: action.confirmation,
    };

    this.isModalOpen = true;
  }

  showTitle() {
    if (this.view) {
      return this.view.totalItems === undefined || this.view.totalItems > 0;
    }
    return true;
  }

  cancelModal() {
    this.resetModal();
  }

  acceptModal() {
    if (this.actionDialogOptions === undefined) {
      return;
    }

    const action = this.actionDialogOptions.action;
    const actionPath = this.actionDialogOptions.action.actionPath;
    const update = { ...action.payload, action: actionPath };
    this.actionService.perform(update);

    this.resetModal();
  }

  private resetModal() {
    this.isModalOpen = false;
    this.actionDialogOptions = undefined;
  }
}

interface ActionDialogOptions {
  action: GridAction;
  text: string;
  type: string;
  confirmation?: Confirmation;
}
