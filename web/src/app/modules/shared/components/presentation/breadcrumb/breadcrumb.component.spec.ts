// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';
import { SharedModule } from '../../../shared.module';
import { BreadcrumbComponent } from './breadcrumb.component';
import { By } from '@angular/platform-browser';
import { DebugElement } from '@angular/core';
import { OverlayScrollbarsComponent } from 'overlayscrollbars-ngx';

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
    component.path = [{ title: 'breadcrumb title', url: '' }];
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

    expect(breadcrumbElement.children.length).toEqual(2);
    expect(breadcrumbElement.innerText).toBe('');
  });

  it('should create two paths with single link', () => {
    component.path = [
      { title: 'breadcrumb title', url: 'some-url' },
      { title: '2nd title', url: '' },
    ];
    const breadcrumbElement: HTMLDivElement = fixture.debugElement.query(
      By.css('.breadcrumb')
    ).nativeElement;
    fixture.detectChanges();

    const links: DebugElement[] = fixture.debugElement.queryAll(By.css('a'));
    const icons: DebugElement[] = fixture.debugElement.queryAll(
      By.css('.separator')
    );
    const spans: DebugElement[] = fixture.debugElement.queryAll(By.css('span'));
    expect(links.length).toBe(1);
    expect(spans.length).toBe(1);

    expect(breadcrumbElement.children.length).toEqual(2);
    expect(breadcrumbElement.innerText).toBe('breadcrumb title \n2nd title');
  });
});
