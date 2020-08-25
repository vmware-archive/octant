// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';

import { GraphvizView } from 'src/app/modules/shared/models/content';
import { D3GraphvizService } from '../../../services/d3/d3graphviz.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-graphviz',
  template: ` <div class="graphviz" #viewer></div> `,
  styleUrls: ['./graphviz.component.scss'],
  encapsulation: ViewEncapsulation.None,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GraphvizComponent extends AbstractViewComponent<GraphvizView> {
  @ViewChild('viewer', { static: true }) private viewer: ElementRef;

  constructor(private d3GraphvizService: D3GraphvizService) {
    super();
  }

  update() {
    const current = this.v.config.dot;
    if (current) {
      this.d3GraphvizService.render(this.viewer, current);
    }
  }
}
