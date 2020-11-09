// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { DropdownComponent } from './dropdown.component';
import { SharedModule } from '../../../shared.module';
import { View, DropdownView } from '../../../models/content';
import { By } from '@angular/platform-browser';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import { instance, mock } from 'ts-mockito';

describe('DropdownComponent', () => {
  let component: DropdownComponent;
  let fixture: ComponentFixture<DropdownComponent>;
  const mockWebsocketService: WebsocketService = mock(WebsocketService);

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [DropdownComponent],
        providers: [
          {
            provide: WebsocketService,
            useValue: instance(mockWebsocketService),
          },
        ],
        imports: [SharedModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(DropdownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    const root: HTMLElement = fixture.nativeElement;
    const el: SVGPathElement = root.querySelector('.dropdown');

    expect(component).toBeTruthy();
    expect(el).not.toBeNull();
  });

  it('should be opened the button dropdown', () => {
    const root: HTMLElement = fixture.nativeElement;
    const view: DropdownView = {
      config: {
        position: 'bottom-left',
        type: 'button',
        action: 'action',
        useSelection: false,
        items: [
          {
            name: 'item-header',
            type: 'header',
            label: 'List of items',
          },
          {
            name: 'first',
            type: 'text',
            label: 'First Item',
          },
          {
            name: 'second',
            type: 'link',
            label: 'Second Item',
            url: '/items/second',
          },
          {
            name: 'item-separator',
            type: 'separator',
          },
          {
            name: 'third',
            type: 'text',
            label: 'Very Last Item',
          },
        ],
      },
      metadata: {
        type: 'dropdown',
        title: [
          {
            metadata: { type: 'text' },
            config: { value: 'Select any item' },
          } as View,
        ],
      },
    };
    component.view = view;

    fixture.detectChanges();
    const dropdown: SVGPathElement = root.querySelector(
      '.dropdown clr-dropdown'
    );
    const dropdownButton: SVGPathElement = root.querySelector(
      '.dropdown clr-dropdown button'
    );
    expect(dropdownButton).not.toBeNull();
    expect(component.position).toBe(view.config.position);

    const dropDownToggle = fixture.debugElement.query(
      By.css('button.dropdown-toggle')
    ).nativeElement;
    dropDownToggle.click();
    fixture.detectChanges();

    expect(dropdown.classList.contains('open')).toBeTruthy();
  });

  it('handles useSelection properly', () => {
    const root: HTMLElement = fixture.nativeElement;
    const view: DropdownView = {
      config: {
        position: 'bottom-left',
        type: 'button',
        action: 'action',
        useSelection: false,
        items: [
          {
            name: 'first',
            type: 'text',
            label: 'First Item',
          },
          {
            name: 'second',
            type: 'link',
            label: 'Second Item',
            url: '/items/second',
          },
        ],
      },
      metadata: {
        type: 'dropdown',
        title: [
          {
            metadata: { type: 'text' },
            config: { value: 'Select any item' },
          } as View,
        ],
      },
    };
    component.view = view;

    fixture.detectChanges();
    const dropdown: SVGPathElement = root.querySelector(
      '.dropdown clr-dropdown'
    );
    const dropDownToggle = fixture.debugElement.query(
      By.css('button.dropdown-toggle')
    ).nativeElement;

    expect(dropDownToggle).not.toBeNull();
    dropDownToggle.click();
    fixture.detectChanges();

    expect(dropdown.classList.contains('open')).toBeTruthy();
    fixture.debugElement
      .query(By.css('clr-dropdown-menu :first-child'))
      .nativeElement.click();

    fixture.detectChanges();
    expect(dropDownToggle.innerHTML).toContain('Select any item');

    view.config.useSelection = true;
    component.view = view;
    dropDownToggle.click();
    fixture.detectChanges();

    fixture.debugElement
      .query(By.css('clr-dropdown-menu :first-child'))
      .nativeElement.click();
    fixture.detectChanges();
    expect(dropDownToggle.innerHTML).toContain('First Item');
  });
});
