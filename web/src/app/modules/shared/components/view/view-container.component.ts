/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ComponentFactoryResolver,
  ComponentRef,
  EventEmitter,
  Inject,
  Input,
  OnInit,
  Output,
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
export class ViewContainerComponent implements OnInit, AfterViewInit {
  @ViewChild(ViewHostDirective, { static: true }) appView: ViewHostDirective;
  lastViewType: string;

  @Input() set view(v: View) {
    if (v && v.metadata) {
      const cur = JSON.stringify(v);
      if (this.previous !== cur) {
        this.previous = cur;
        this.loadView(v);
      }
    }
  }
  @Input() enableDebug = false;
  @Output() viewInit: EventEmitter<void> = new EventEmitter<void>();

  private start: number;
  private componentRef: ComponentRef<Viewer>;
  private previous: string;

  constructor(
    private componentFactoryResolver: ComponentFactoryResolver,
    @Inject(DYNAMIC_COMPONENTS_MAPPING)
    private componentMappings: ComponentMapping
  ) {}

  ngOnInit(): void {
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

  loadView(view: View) {
    if (!this.componentRef || this.lastViewType !== view.metadata.type) {
      const viewType = view.metadata.type;
      let component: Type<any> = this.componentMappings[viewType];
      if (!component) {
        component = MissingComponentComponent;
      }

      const componentFactory = this.componentFactoryResolver.resolveComponentFactory(
        component
      );
      const viewContainerRef = this.appView.viewContainerRef;
      viewContainerRef.clear();

      this.componentRef = viewContainerRef.createComponent<Viewer>(
        componentFactory
      );
    }
    this.componentRef.instance.view = view;
    this.componentRef.instance.viewInit.subscribe(_ => this.viewInit.emit());
  }
}
