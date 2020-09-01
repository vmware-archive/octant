// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  EventEmitter,
  OnDestroy,
  Output,
  ViewEncapsulation,
} from '@angular/core';
import { PortsView } from 'src/app/modules/shared/models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-ports',
  templateUrl: './ports.component.html',
  styleUrls: ['./ports.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class PortsComponent
  extends AbstractViewComponent<PortsView>
  implements OnDestroy {
  private previousView: PortsView;

  @Output() portLoad: EventEmitter<boolean> = new EventEmitter(true);

  constructor() {
    super();
  }

  update() {
    if (JSON.stringify(this.v) !== JSON.stringify(this.previousView)) {
      this.portLoad.emit(false);
    }
  }

  ngOnDestroy() {}

  load(_: Event) {
    this.portLoad.emit(true);
  }

  trackByFn(index, _) {
    return index;
  }
}
