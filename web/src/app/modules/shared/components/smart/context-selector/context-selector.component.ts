// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import {
  ContextDescription,
  KubeContextService,
} from '../../../services/kube-context/kube-context.service';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-context-selector',
  templateUrl: './context-selector.component.html',
  styleUrls: ['./context-selector.component.scss'],
})
export class ContextSelectorComponent implements OnInit, OnDestroy {
  contexts: ContextDescription[];
  selected: string;

  private kubeContextSubscription: Subscription;

  constructor(private kubeContext: KubeContextService) {}

  ngOnInit() {
    this.kubeContextSubscription = this.kubeContext
      .contexts()
      .subscribe(contexts => (this.contexts = contexts));
    this.kubeContextSubscription = this.kubeContext
      .selected()
      .subscribe(selected => (this.selected = selected));
  }

  ngOnDestroy(): void {
    if (this.kubeContextSubscription) {
      this.kubeContextSubscription.unsubscribe();
    }
  }

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
