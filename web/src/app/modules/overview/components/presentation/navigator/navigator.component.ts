/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
} from '@angular/core';
import { NavigatorService } from '../../../services/navigator/navigator.service';

@Component({
  selector: 'app-navigator',
  templateUrl: './navigator.component.html',
  styleUrls: ['./navigator.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NavigatorComponent implements OnInit {
  backLink: string;

  private index: number;
  private history: string[];

  constructor(
    private navigator: NavigatorService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.navigator.status.subscribe(status => {
      this.index = status.index;
      this.history = status.history;

      if (this.index > 0) {
        this.backLink = this.history[this.index - 1];
        this.cdr.detectChanges();
      } else {
        this.backLink = '';
      }
    });
  }

  backStyle() {
    if (this.backLink.length < 1) {
      return {
        display: 'none',
      };
    }

    return {};
  }
}
