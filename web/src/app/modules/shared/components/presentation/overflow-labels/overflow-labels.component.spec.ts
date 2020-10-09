// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { OverflowLabelsComponent } from './overflow-labels.component';
import { By } from '@angular/platform-browser';
import { LabelFilterService } from '../../../services/label-filter/label-filter.service';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import { WebsocketServiceMock } from '../../../../../data/services/websocket/mock';
import { windowProvider, WindowToken } from '../../../../../window';

describe('OverflowLabelsComponent', () => {
  let component: OverflowLabelsComponent;
  let fixture: ComponentFixture<OverflowLabelsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [OverflowLabelsComponent],
      providers: [
        LabelFilterService,
        { provide: WindowToken, useFactory: windowProvider },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OverflowLabelsComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display only two labels with a +1 badge', () => {
    component.labels = {
      ['keyOne']: 'valueOne',
      ['keyTwo']: 'valueTwo',
      ['keyThree']: 'valueThree',
    };
    fixture.detectChanges();
    const renderedLabels = document.getElementsByClassName('label');

    expect(component.overflowLabels.length).toEqual(1);
    expect(component.showLabels.length).toEqual(2);
    expect(renderedLabels.length).toEqual(2);
  });

  it('should display all labels if the number is less or equal than the number to display', () => {
    component.labels = {
      ['keyOne']: 'valueOne',
    };
    fixture.detectChanges();
    const renderedLabels = document.getElementsByClassName('label');

    expect(component.overflowLabels).toBeUndefined();
    expect(component.showLabels.length).toEqual(1);
    expect(renderedLabels.length).toEqual(1);
  });

  it('should call addFilter method when clicking on a label', () => {
    spyOn(component, 'filterLabel');
    component.labels = {
      ['keyOne']: 'valueOne',
    };
    fixture.detectChanges();
    const firstLabel = fixture.debugElement.query(By.css('.label'))
      .nativeElement;

    firstLabel.click();
    expect(component.filterLabel).toHaveBeenCalledWith('keyOne', 'valueOne');
  });

  it('should add the correct filter', () => {
    const debugElement = fixture.debugElement;
    const labelFilterService = debugElement.injector.get(LabelFilterService);
    spyOn(labelFilterService, 'add');
    component.filterLabel('keyOne', 'valueOne');

    expect(labelFilterService.add).toHaveBeenCalledWith({
      key: 'keyOne',
      value: 'valueOne',
    });
  });
});
