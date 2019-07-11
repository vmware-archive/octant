// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { Action, SummaryItem, SummaryView } from 'src/app/models/content';
import { FormGroup } from '@angular/forms';
import { ActionService } from '../../services/action/action.service';
import { ViewService } from '../../services/view/view.service';

@Component({
  selector: 'app-view-summary',
  templateUrl: './summary.component.html',
  styleUrls: ['./summary.component.scss'],
})
export class SummaryComponent implements OnChanges {
  @Input() view: SummaryView;
  title: string;
  isLoading = false;

  currentAction: Action;

  constructor(
    private actionService: ActionService,
    private viewService: ViewService
  ) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view: SummaryView = changes.view.currentValue;
      this.title = this.viewService.viewTitleAsText(view);
    }
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
      this.actionService.perform(formGroup.value).subscribe();
      this.currentAction = undefined;
    }
  }

  onActionCancel() {
    this.currentAction = undefined;
  }

  shouldShowFooter(): boolean {
    if (this.view && this.view.config.actions) {
      if (!this.currentAction && this.view.config.actions.length > 0) {
        return true;
      }
    }

    return false;
  }
}
