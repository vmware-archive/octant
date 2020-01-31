// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Location } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { Injectable, NgModule, NgZone } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { Router, RouterModule } from '@angular/router';
import { ClarityModule } from '@clr/angular';
import { NgSelectModule } from '@ng-select/ng-select';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { InputFilterComponent } from './components/input-filter/input-filter.component';
import { NamespaceComponent } from './components/namespace/namespace.component';
import { NavigationComponent } from './components/navigation/navigation.component';
import { QuickSwitcherComponent } from './components/quick-switcher/quick-switcher.component';
import { NotifierComponent } from './components/notifier/notifier.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { OverviewModule } from './modules/overview/overview.module';
import { MarkdownModule, MarkedOptions } from 'ngx-markdown';
import { ThemeSwitchButtonComponent } from './modules/overview/components/theme-switch/theme-switch-button.component';
import { SharedModule } from './shared/shared.module';

@Injectable()
export class UnstripTrailingSlashLocation extends Location {
  public static stripTrailingSlash(url: string): string {
    return url;
  }
}

@NgModule({
  declarations: [
    AppComponent,
    NamespaceComponent,
    PageNotFoundComponent,
    InputFilterComponent,
    NotifierComponent,
    NavigationComponent,
    QuickSwitcherComponent,
    ThemeSwitchButtonComponent,
  ],
  imports: [
    BrowserModule,
    ClarityModule,
    BrowserAnimationsModule,
    HttpClientModule,
    RouterModule,
    FormsModule,
    AppRoutingModule,
    OverviewModule,
    NgSelectModule,
    MarkdownModule.forRoot({
      markedOptions: {
        provide: MarkedOptions,
        useValue: {
          gfm: true,
          tables: true,
          breaks: true,
          pedantic: false,
          sanitize: false,
          smartLists: true,
          smartypants: false,
        },
      },
    }),
    SharedModule,
  ],
  providers: [
    {
      provide: Location,
      useClass: UnstripTrailingSlashLocation,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
  constructor(private ngZone: NgZone, private router: Router) {}

  navigate(commands: any[]): void {
    this.ngZone.run(() => this.router.navigate(commands)).then();
  }
}
