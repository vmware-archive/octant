// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component } from '@angular/core';
import {
  Action,
  SummaryItem,
  SummaryView,
} from 'src/app/modules/shared/models/content';
import { FormGroup } from '@angular/forms';
import { ActionService } from '../../../services/action/action.service';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-summary',
  templateUrl: './summary.component.html',
  styleUrls: ['./summary.component.scss'],
})
export class SummaryComponent extends AbstractViewComponent<SummaryView> {
  title: string;
  isLoading = false;

  currentAction: Action;

  constructor(
    private actionService: ActionService,
    private viewService: ViewService
  ) {
    super();
  }

  update() {
    const view = this.v;
    this.title = this.viewService.viewTitleAsText(view);
  }

  identifyItem(index: number, item: SummaryItem): string {
    return `${index}-${item.header}`;
  }

  onPortLoad(isLoading: boolean) {
    this.isLoading = isLoading;
  }

  setAction(action: Action) {
    this.currentAction = action;
  }

  onActionSubmit(formGroup: FormGroup) {
    if (formGroup && formGroup.value) {
      this.actionService.perform(formGroup.value);
      this.currentAction = undefined;
    }
  }

  onActionCancel() {
    this.currentAction = undefined;
  }

  shouldShowFooter(): boolean {
    if (this.v && this.v.config.actions) {
      if (!this.currentAction && this.v.config.actions.length > 0) {
        return true;
      }
    }

    return false;
  }
}
