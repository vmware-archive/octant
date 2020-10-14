// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SharedModule } from '../../../shared.module';
import { OctantTooltipComponent } from '../octant-tooltip/octant-tooltip';
import { OverflowSelectorsComponent } from './overflow-selectors.component';
import { windowProvider, WindowToken } from '../../../../../window';

describe('OverflowSelectorsComponent', () => {
  let component: OverflowSelectorsComponent;
  let fixture: ComponentFixture<OverflowSelectorsComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [OverflowSelectorsComponent, OctantTooltipComponent],
        imports: [SharedModule],
        providers: [{ provide: WindowToken, useFactory: windowProvider }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(OverflowSelectorsComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display only two selectors with a +1 badge', () => {
    component.selectors = [
      {
        metadata: {
          type: 'LabelSelector',
        },
        config: {
          key: 'keyOne',
          value: 'valueOne',
        },
      },
      {
        metadata: {
          type: 'LabelSelector',
        },
        config: {
          key: 'keyTwo',
          value: 'valueTwo',
        },
      },
      {
        metadata: {
          type: 'LabelSelector',
        },
        config: {
          key: 'keyThree',
          value: 'valueThree',
        },
      },
    ];

    fixture.detectChanges();
    const renderedSelectors = document.getElementsByClassName('label');

    expect(component.overflowSelectors.length).toEqual(1);
    expect(component.showSelectors.length).toEqual(2);
    expect(renderedSelectors.length).toEqual(2);
  });

  it('should display all selectors if the number is less or equal than the number to display', () => {
    component.selectors = [
      {
        metadata: {
          type: 'LabelSelector',
        },
        config: {
          key: 'keyOne',
          value: 'valueOne',
        },
      },
    ];
    fixture.detectChanges();
    const renderedSelectors = document.getElementsByClassName('label');

    expect(component.overflowSelectors).toBeUndefined();
    expect(component.showSelectors.length).toEqual(1);
    expect(renderedSelectors.length).toEqual(1);
  });
});
