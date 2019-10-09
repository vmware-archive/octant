// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ContextSelectorComponent } from './context-selector.component';
import { KubeContextService } from '../../services/kube-context/kube-context.service';
import { of } from 'rxjs';
import { By } from '@angular/platform-browser';

const contexts = [
  {
    name: 'kubernetes-admin@service-account',
  },
  {
    name: 'kubernetes-admin@workload-test',
  },
];

class MockKubeContextService {
  contexts() {
    return of(contexts);
  }

  selected() {
    return of(contexts[0].name);
  }
}

describe('ContextSelectorComponent', () => {
  let component: ContextSelectorComponent;
  let fixture: ComponentFixture<ContextSelectorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ContextSelectorComponent],
      providers: [
        { provide: KubeContextService, useClass: MockKubeContextService },
      ],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContextSelectorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('selects the context when the button is clicked', () => {
    const onClickMock = spyOn(component, 'selectContext');

    const dropDownToggle = fixture.debugElement.query(
      By.css('button.dropdown-toggle')
    ).nativeElement;
    dropDownToggle.click();

    fixture.detectChanges();

    fixture.debugElement
      .query(By.css('button.context-button:last-child'))
      .nativeElement.click();

    fixture.detectChanges();

    expect(onClickMock).toHaveBeenCalledWith(contexts[1]);
  });

  it('makes the currently selected context active', () => {
    component.selected = contexts[1].name;

    const dropDownToggle = fixture.debugElement.query(
      By.css('button.dropdown-toggle')
    ).nativeElement;
    dropDownToggle.click();

    fixture.detectChanges();

    const contextButton = fixture.debugElement.query(
      By.css('button.context-button:last-child')
    );

    expect(contextButton.classes.active).toBeTruthy();
  });

  it('shows the currently active context name', () => {
    const dropDownToggle: HTMLButtonElement = fixture.debugElement.query(
      By.css('button.dropdown-toggle')
    ).nativeElement;

    expect(dropDownToggle.textContent.trim()).toBe(contexts[0].name);
  });
});
