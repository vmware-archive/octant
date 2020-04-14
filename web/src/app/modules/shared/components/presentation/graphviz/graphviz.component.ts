// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  AfterViewChecked,
  Component,
  ElementRef,
  Input,
  ViewChild,
  ViewEncapsulation,
} from '@angular/core';

import { GraphvizView, View } from '../../../../shared/models/content';
import { D3GraphvizService } from '../../../services/d3/d3graphviz.service';

@Component({
  selector: 'app-view-graphviz',
  template: ` <div class="graphviz" #viewer></div> `,
  styleUrls: ['./graphviz.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class GraphvizComponent implements AfterViewChecked {
  @ViewChild('viewer', { static: true }) private viewer: ElementRef;

  private v: GraphvizView;

  @Input() set view(v: View) {
    this.v = v as GraphvizView;
  }
  get view() {
    return this.v;
  }

  constructor(private d3GraphvizService: D3GraphvizService) {}

  ngAfterViewChecked() {
    if (this.view) {
      const current = this.v.config.dot;
      this.d3GraphvizService.render(this.viewer, current);
    }
  }
}
