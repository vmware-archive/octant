// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule } from '@angular/forms';
import { InputFilterComponent } from './input-filter.component';
import {
  LabelFilterService,
  Filter,
} from 'src/app/shared/services/label-filter/label-filter.service';
import { BehaviorSubject } from 'rxjs';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';

const labelFilterStub: Partial<LabelFilterService> = {
  filters: new BehaviorSubject<Filter[]>([]),
};

describe('InputFilterComponent', () => {
  let component: InputFilterComponent;
  let fixture: ComponentFixture<InputFilterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [FormsModule],
      declarations: [InputFilterComponent],
      providers: [{ provide: LabelFilterService, useValue: labelFilterStub }],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InputFilterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeDefined();
  });

  it('should not show the tag list on init render', () => {
    const tagListDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.input-filter-tags')
    );
    expect(tagListDebugElement).toBeNull();
  });

  it('should show the tag list if down arrow icon clicked', () => {
    expect(component.showTagList).toBe(false);
    const downArrowIconDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.down-icon')
    );
    downArrowIconDebugElement.triggerEventHandler('click', null);
    expect(component.showTagList).toBe(true);
    fixture.detectChanges();
    const tagListElement: HTMLElement = fixture.debugElement.query(
      By.css('.input-filter-tags')
    ).nativeElement;
    expect(tagListElement).not.toBeNull();
  });

  it('should show the user text if there are no filters', () => {
    component.showTagList = true;
    fixture.detectChanges();
    const fixtureDebugElement: DebugElement = fixture.debugElement;
    const userTextDebugElement: DebugElement = fixtureDebugElement.query(
      By.css('.input-filter-empty')
    );
    const userTextNativeElement: HTMLElement =
      userTextDebugElement.nativeElement;
    expect(userTextNativeElement.textContent).toMatch(/No current filters/i);
  });

  it('should show the tags if there are filters', () => {
    const labelFilterService: LabelFilterService = TestBed.get(
      LabelFilterService
    );
    labelFilterService.filters.next([
      { key: 'test1', value: 'filter1' },
      { key: 'test2', value: 'filter2' },
      { key: 'test3', value: 'filter3' },
    ]);
    component.showTagList = true;
    fixture.detectChanges();
    const tagDebugElements: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.input-filter-tag')
    );
    expect(tagDebugElements.length).toBe(3);
    expect(tagDebugElements[0].nativeElement.textContent).toMatch(
      /test1:filter1/i
    );
    expect(tagDebugElements[1].nativeElement.textContent).toMatch(
      /test2:filter2/i
    );
    expect(tagDebugElements[2].nativeElement.textContent).toMatch(
      /test3:filter3/i
    );
  });

  it('should change the placeholder text if filters are applied', () => {
    let inputElement: HTMLInputElement = fixture.debugElement.query(
      By.css('.text-input')
    ).nativeElement;
    expect(inputElement.placeholder).toMatch(/\bFilter by labels\b/i);

    const labelFilterService: LabelFilterService = TestBed.get(
      LabelFilterService
    );
    labelFilterService.filters.next([
      { key: 'test1', value: 'filter1' },
      { key: 'test2', value: 'filter2' },
      { key: 'test3', value: 'filter3' },
    ]);
    fixture.detectChanges();
    inputElement = fixture.debugElement.query(By.css('.text-input'))
      .nativeElement;
    expect(inputElement.placeholder).toMatch(/Filter by labels \(3 applied\)/i);

    labelFilterService.filters.next([{ key: 'test1', value: 'filter1' }]);
    fixture.detectChanges();
    inputElement = fixture.debugElement.query(By.css('.text-input'))
      .nativeElement;
    expect(inputElement.placeholder).toMatch(/Filter by labels \(1 applied\)/i);
  });

  it('should be able to enter a tag through the input', async(() => {
    const inputDebugElement: DebugElement = fixture.debugElement.query(
      By.css('.text-input')
    );
    let inputNativeElement: HTMLInputElement = inputDebugElement.nativeElement;
    inputNativeElement.value = 'test1:filter1';
    inputNativeElement.dispatchEvent(new Event('input'));
    expect(component.inputValue).toBe('test1:filter1');

    const labelFilterService = TestBed.get(LabelFilterService);
    labelFilterService.decodeFilter = jasmine
      .createSpy('decodeFilter')
      .and.returnValue({ key: 'test1', value: 'filter1' });
    labelFilterService.add = jasmine.createSpy('add');
    inputDebugElement.triggerEventHandler('keyup.enter', null);
    expect(labelFilterService.decodeFilter.calls.count()).toBe(1);
    expect(labelFilterService.add.calls.count()).toBe(1);
    expect(component.showTagList).toBe(true);
    expect(component.inputValue).toBe('');

    fixture.detectChanges();
    fixture.whenStable().then(() => {
      inputNativeElement = fixture.debugElement.query(By.css('.text-input'))
        .nativeElement;
      expect(inputNativeElement.value).toBe('');
    });
  }));
});
