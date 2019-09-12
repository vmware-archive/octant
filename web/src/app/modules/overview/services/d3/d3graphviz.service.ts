// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ElementRef, Injectable } from '@angular/core';
import { graphviz } from 'd3-graphviz';

@Injectable({
  providedIn: 'root',
})
export class D3GraphvizService {
  constructor() {}

  render(parentElement: ElementRef, g) {
    const viewer = parentElement.nativeElement;
    graphviz(viewer).renderDot(g);
  }
}
