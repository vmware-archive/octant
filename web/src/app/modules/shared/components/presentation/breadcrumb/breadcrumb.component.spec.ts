// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SharedModule } from '../../../shared.module';
import { BreadcrumbComponent } from './breadcrumb.component';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { OverlayScrollbarsComponent } from 'overlayscrollbars-ngx';
import { LinkView, TextView } from '../../../models/content';

describe('BreadcrumbComponent', () => {
  let component: BreadcrumbComponent;
  let fixture: ComponentFixture<BreadcrumbComponent>;

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [OverlayScrollbarsComponent],
        imports: [SharedModule],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(BreadcrumbComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should omit if single path', () => {
    const text: TextView = {
      config: { value: 'breadcrumb title' },
      metadata: { type: 'text' },
    };
    component.path = [text];
    const breadcrumbElement: HTMLDivElement = fixture.debugElement.query(
      By.css('.breadcrumb')
    ).nativeElement;
    fixture.detectChanges();

    const links: DebugElement[] = fixture.debugElement.queryAll(By.css('a'));
    const icons: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.separator')
    );
    const spans: DebugElement[] = fixture.debugElement.queryAll(By.css('span'));
    expect(links.length).toBe(0);
    expect(spans.length).toBe(0);
    expect(icons.length).toBe(0);

    expect(breadcrumbElement.children.length).toEqual(1);
    expect(breadcrumbElement.innerText).toBe('');
  });

  it('should create two paths with single link', () => {
    const link: LinkView = {
      config: { value: 'breadcrumb title', ref: 'some-url' },
      metadata: { type: 'link' },
    };
    const text: TextView = {
      config: { value: '2nd title' },
      metadata: { type: 'text' },
    };
    component.path = [link, text];
    const breadcrumbElement: HTMLDivElement = fixture.debugElement.query(
      By.css('.breadcrumb')
    ).nativeElement;
    fixture.detectChanges();

    const links: DebugElement[] = fixture.debugElement.queryAll(By.css('a'));
    const icons: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.separator')
    );
    const spans: DebugElement[] = fixture.debugElement.queryAll(
      By.css('app-view-text')
    );
    expect(links.length).toBe(1);
    expect(spans.length).toBe(1);
    expect(icons.length).toBe(1);

    expect(breadcrumbElement.children.length).toEqual(2);
    expect(breadcrumbElement.innerText).toContain('breadcrumb title');
    expect(breadcrumbElement.innerText).toContain('2nd title');
  });
});
