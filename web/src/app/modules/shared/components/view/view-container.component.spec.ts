/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ViewContainerComponent } from './view-container.component';
import { DYNAMIC_COMPONENTS_MAPPING } from '../../dynamic-components';
import { TextComponent } from '../presentation/text/text.component';
import { PodStatusView, TextView } from '../../models/content';
import { SharedModule } from '../../shared.module';
import { PodStatusComponent } from '../presentation/pod-status/pod-status.component';

describe('ViewContainerComponent', () => {
  let component: ViewContainerComponent;
  let fixture: ComponentFixture<ViewContainerComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ViewContainerComponent],
        imports: [SharedModule],
        providers: [
          {
            provide: DYNAMIC_COMPONENTS_MAPPING,
            useValue: {
              text: TextComponent,
              podStatus: PodStatusComponent,
            },
          },
        ],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ViewContainerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should recreate component when different type', () => {
    const textView: TextView = {
      config: { value: 'some text' },
      metadata: { type: 'text', title: [], accessor: 'accessor' },
    };

    const podStatusView: PodStatusView = {
      metadata: { type: 'podStatus' },
      config: {
        pods: {
          pod1: {
            details: [textView],
            status: 'ok',
          },
        },
      },
    };

    component.view = textView;
    fixture.detectChanges();
    expect(component.componentRef.instance.view.metadata.type).toEqual('text');
    expect(component.componentRef.componentType.name).toEqual('TextComponent');

    component.view = podStatusView;
    fixture.detectChanges();
    expect(component.componentRef.instance.view.metadata.type).toEqual(
      'podStatus'
    );
    expect(component.componentRef.componentType.name).toEqual(
      'PodStatusComponent'
    );
  });
});
