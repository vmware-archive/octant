// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { ExpressionSelectorComponent } from './expression-selector.component';
import { ExpressionSelectorView } from 'src/app/modules/shared/models/content';

describe('ExpressionSelectorComponent', () => {
  let component: ExpressionSelectorComponent;
  let fixture: ComponentFixture<ExpressionSelectorComponent>;
  let view: ExpressionSelectorView;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [ExpressionSelectorComponent],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(ExpressionSelectorComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();

    expect(component).toBeTruthy();
  });

  describe('expression with single value', () => {
    beforeEach(() => {
      view = {
        metadata: {
          type: 'expressionSelector',
        },
        config: {
          key: 'keyOne',
          operator: 'NotIn',
          values: ['valueOne'],
        },
      };

      component.view = view;
      fixture.detectChanges();
    });

    it('should render correctly', () => {
      const element: HTMLElement = fixture.nativeElement.querySelector('div');
      expect(element.textContent).toEqual('keyOne NotIn valueOne');
    });
  });

  describe('expression with multiple values', () => {
    beforeEach(() => {
      view = {
        metadata: {
          type: 'expressionSelector',
        },
        config: {
          key: 'keyOne',
          operator: 'In',
          values: ['valueOne', 'valueTwo'],
        },
      };

      component.view = view;
      fixture.detectChanges();
    });

    it('should render correctly', () => {
      const element: HTMLElement = fixture.nativeElement.querySelector('div');
      expect(element.textContent).toEqual('keyOne In valueOne|valueTwo');
    });
  });

  describe('expression with no values', () => {
    beforeEach(() => {
      view = {
        metadata: {
          type: 'expressionSelector',
        },
        config: {
          key: 'keyOne',
          operator: 'Exists',
          values: [],
        },
      };

      component.view = view;
      fixture.detectChanges();
    });

    it('should render correctly', () => {
      const element: HTMLElement = fixture.nativeElement.querySelector('div');
      expect(element.textContent).toEqual('keyOne Exists');
    });
  });
});
