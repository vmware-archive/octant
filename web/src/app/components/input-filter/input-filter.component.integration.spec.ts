// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { LabelFilterService, Filter } from 'src/app/services/label-filter/label-filter.service';
import { InputFilterComponent } from './input-filter.component';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

describe('InputFilterComponent - LabelFilterService integration test', () => {
  let component: InputFilterComponent;
  let fixture: ComponentFixture<InputFilterComponent>;
  let service: LabelFilterService;
  let router: Router;

  beforeEach(async(() => {
    const routerStub = {
      events: { subscribe: jasmine.createSpy() },
      routerState: { root: { queryParamMap: { subscribe: jasmine.createSpy() }}},
      navigate: jasmine.createSpy(),
    };

    TestBed.configureTestingModule({
      declarations: [
        InputFilterComponent,
      ],
      providers: [
        { provide: Router, useValue: routerStub },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    service = TestBed.get(LabelFilterService);
    router = TestBed.get(Router);
    fixture = TestBed.createComponent(InputFilterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeDefined();
    expect(service).toBeDefined();
    const filters = service.filters.getValue();
    expect(filters).toEqual([]);
  });

  it('should update placeholder text when service adds filter', () => {
    let inputElement: HTMLInputElement = fixture.debugElement.query(By.css('.text-input')).nativeElement;
    expect(inputElement.placeholder).toMatch(/\bFilter by labels\b/i);

    const testFilter = { key: 'test1', value: 'value1' };
    service.add(testFilter);
    expect((router.navigate as jasmine.Spy).calls.count()).toBe(1);

    fixture.detectChanges();
    inputElement = fixture.debugElement.query(By.css('.text-input')).nativeElement;
    expect(inputElement.placeholder).toMatch(/Filter by labels \(1 applied\)/i);

    const testFilter2 = { key: 'test2', value: 'value2' };
    service.add(testFilter2);
    fixture.detectChanges();
    inputElement = fixture.debugElement.query(By.css('.text-input')).nativeElement;
    expect(inputElement.placeholder).toMatch(/Filter by labels \(2 applied\)/i);
  });

  it('should update service if filter is removed from ui', () => {
    service.add({ key: 'test1', value: 'value1' });
    service.add({ key: 'test2', value: 'value2' });
    service.add({ key: 'test3', value: 'value3' });
    component.showTagList = true;
    fixture.detectChanges();

    const tagDebugElements: DebugElement[] = fixture.debugElement.queryAll(By.css('.input-filter-tag'));
    expect(tagDebugElements.length).toBe(3);
    expect(tagDebugElements[0].nativeElement.textContent).toMatch(/test1:value1/i);
    expect(tagDebugElements[1].nativeElement.textContent).toMatch(/test2:value2/i);
    expect(tagDebugElements[2].nativeElement.textContent).toMatch(/test3:value3/i);

    tagDebugElements[1].query(By.css('.input-filter-tag-remove')).triggerEventHandler('click', null);
    const observedFilters = service.filters.getValue();
    const expectedFilters = [{ key: 'test1', value: 'value1' }, { key: 'test3', value: 'value3' }];
    expect(observedFilters).toEqual(expectedFilters);
    expect((router.navigate as jasmine.Spy).calls.count()).toBe(4);
  });

  it('should update service if filter is added from input', () => {
    const inputDebugElement: DebugElement = fixture.debugElement.query(By.css('.text-input'));
    const inputElement: HTMLInputElement = inputDebugElement.nativeElement;
    inputElement.value = 'test1:value1';
    inputElement.dispatchEvent(new Event('input'));
    inputDebugElement.triggerEventHandler('keyup.enter', null);

    const observedFilters = service.filters.getValue();
    const expectedFilters: Filter[] = [{ key: 'test1', value: 'value1' }];
    expect(observedFilters).toEqual(expectedFilters);
  });
});
