// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';

import { GraphvizView, View } from 'src/app/modules/shared/models/content';
import { D3GraphvizService } from '../../../services/d3/d3graphviz.service';

@Component({
  selector: 'app-view-graphviz',
  template: ` <div class="graphviz" #viewer></div> `,
  styleUrls: ['./graphviz.component.scss'],
  encapsulation: ViewEncapsulation.None,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GraphvizComponent implements OnChanges {
  @ViewChild('viewer', { static: true }) private viewer: ElementRef;

  private v: GraphvizView;

  @Input() set view(v: View) {
    this.v = v as GraphvizView;
  }
  get view() {
    return this.v;
  }

  constructor(private d3GraphvizService: D3GraphvizService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as GraphvizView;
      const current = view.config.dot;
      this.d3GraphvizService.render(this.viewer, current);
    }
  }
}
