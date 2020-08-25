/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ComponentFactoryResolver,
  EventEmitter,
  Inject,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
  Type,
  ViewChild,
} from '@angular/core';
import { View } from '../../models/content';
import { ViewHostDirective } from '../../directives/view-host/view-host.directive';
import {
  ComponentMapping,
  DYNAMIC_COMPONENTS_MAPPING,
} from '../../dynamic-components';
import { MissingComponentComponent } from '../missing-component/missing-component.component';

interface Viewer {
  view: View;
  viewInit: EventEmitter<void>;
  ping: () => void;
}

@Component({
  selector: 'app-view-container',
  template: `<ng-container appView></ng-container>`,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ViewContainerComponent
  implements OnInit, OnChanges, AfterViewInit {
  @ViewChild(ViewHostDirective, { static: true }) appView: ViewHostDirective;
  @Input() view: View;
  @Input() enableDebug = false;
  @Output() viewInit: EventEmitter<void> = new EventEmitter<void>();

  private start: number;

  constructor(
    private componentFactoryResolver: ComponentFactoryResolver,
    @Inject(DYNAMIC_COMPONENTS_MAPPING)
    private componentMappings: ComponentMapping
  ) {}

  ngOnInit(): void {
    this.loadView();

    if (this.enableDebug) {
      this.start = new Date().getTime();
    }
  }

  ngAfterViewInit() {
    if (this.view === null) {
      return;
    }

    if (this.enableDebug) {
      console.log(
        `${this.view.metadata.type}: ${new Date().getTime() - this.start}`
      );
    }
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes.view.currentValue) {
      // TODO: send a checksum so this doesn't need to be calculated here.
      const prev = JSON.stringify(changes.view.previousValue);
      const cur = JSON.stringify(changes.view.currentValue);

      if (cur !== prev) {
        this.loadView();
      }
    }
  }

  loadView() {
    if (!this.view || !this.view.metadata) {
      return;
    }

    const viewType = this.view.metadata.type;
    let component: Type<any> = this.componentMappings[viewType];
    if (!component) {
      component = MissingComponentComponent;
    }

    const componentFactory = this.componentFactoryResolver.resolveComponentFactory(
      component
    );
    const viewContainerRef = this.appView.viewContainerRef;
    viewContainerRef.clear();

    const componentRef = viewContainerRef.createComponent<Viewer>(
      componentFactory
    );
    componentRef.instance.view = this.view;
    componentRef.instance.viewInit.subscribe(_ => this.viewInit.emit());
  }
}
