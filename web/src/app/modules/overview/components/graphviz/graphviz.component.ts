// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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

import { GraphvizView } from 'src/app/models/content';
import { D3GraphvizService } from '../../services/d3/d3graphviz.service';

@Component({
  selector: 'app-view-graphviz',
  template: `
    <div class="graphviz" #viewer></div>
  `,
  styleUrls: ['./graphviz.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class GraphvizComponent implements AfterViewChecked {
  @ViewChild('viewer') private viewer: ElementRef;
  @Input() view: GraphvizView;

  constructor(private d3GraphvizService: D3GraphvizService) {}

  ngAfterViewChecked() {
    if (this.view) {
      const current = this.view.config.dot;
      this.d3GraphvizService.render(this.viewer, current);
    }
  }
}
