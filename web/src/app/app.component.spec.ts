// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { HttpClientTestingModule } from '@angular/common/http/testing';
import { async, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { ClarityModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { AppComponent } from './app.component';
import { NamespaceComponent } from './components/namespace/namespace.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { InputFilterComponent } from './components/input-filter/input-filter.component';
import { NotifierComponent } from './components/notifier/notifier.component';
import { NavigationComponent } from './components/navigation/navigation.component';
import { ContextSelectorComponent } from './modules/overview/components/context-selector/context-selector.component';
import { DefaultPipe } from './modules/overview/pipes/default.pipe';
import { NgSelectModule } from '@ng-select/ng-select';

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        RouterTestingModule,
        ClarityModule,
        HttpClientTestingModule,
        FormsModule,
        NgSelectModule,
      ],
      declarations: [
        AppComponent,
        NamespaceComponent,
        PageNotFoundComponent,
        InputFilterComponent,
        NotifierComponent,
        NavigationComponent,
        ContextSelectorComponent,
        DefaultPipe,
      ],
    }).compileComponents();
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  // it(`should have as title 'dash-frontend'`, () => {
  //   const fixture = TestBed.createComponent(AppComponent);
  //   const app = fixture.debugElement.componentInstance;
  //   expect(app.title).toEqual('dash-frontend');
  // });

  // it('should render title in a h1 tag', () => {
  //   const fixture = TestBed.createComponent(AppComponent);
  //   fixture.detectChanges();
  //   const compiled = fixture.debugElement.nativeElement;
  //   expect(compiled.querySelector('h1').textContent).toContain('Welcome to dash-frontend!');
  // });
});
