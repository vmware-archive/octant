// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnInit, OnDestroy } from '@angular/core';
import { untilDestroyed } from 'ngx-take-until-destroy';
import {
  KubeContextService,
  ContextDescription,
} from '../../../../modules/overview/services/kube-context/kube-context.service';

@Component({
  selector: 'app-context-selector',
  templateUrl: './context-selector.component.html',
  styleUrls: ['./context-selector.component.scss'],
})
export class ContextSelectorComponent implements OnInit, OnDestroy {
  contexts: ContextDescription[];
  selected: string;

  constructor(private kubeContext: KubeContextService) {}

  ngOnInit() {
    this.kubeContext
      .contexts()
      .pipe(untilDestroyed(this))
      .subscribe(contexts => (this.contexts = contexts));
    this.kubeContext
      .selected()
      .pipe(untilDestroyed(this))
      .subscribe(selected => (this.selected = selected));
  }

  ngOnDestroy() {}

  contextClass(context: ContextDescription) {
    const active = this.selected === context.name ? ['active'] : [];
    return ['context-button', ...active];
  }

  selectContext(context: ContextDescription) {
    this.kubeContext.select(context);
  }

  trackByFn(index, item) {
    return index;
  }
}
