// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnDestroy,
  Output,
  SimpleChanges,
  ViewEncapsulation,
} from '@angular/core';
import { PortsView, View } from '../../../../shared/models/content';

@Component({
  selector: 'app-view-ports',
  templateUrl: './ports.component.html',
  styleUrls: ['./ports.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class PortsComponent implements OnChanges, OnDestroy {
  v: PortsView;

  @Input() set view(v: View) {
    this.v = v as PortsView;
  }
  get view() {
    return this.v;
  }

  @Output() portLoad: EventEmitter<boolean> = new EventEmitter(true);

  constructor() {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue && !changes.view.isFirstChange()) {
      const currentView = changes.view.currentValue as PortsView;
      const prevView = changes.view.previousValue as PortsView;
      if (JSON.stringify(currentView) !== JSON.stringify(prevView)) {
        this.portLoad.emit(false);
      }
    }
  }

  ngOnDestroy() {}

  load(e: Event) {
    this.portLoad.emit(true);
  }

  trackByFn(index, item) {
    return index;
  }
}
