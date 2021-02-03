// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { CommonModule, Location } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { Injectable, NgModule } from '@angular/core';
import { RouteReuseStrategy, RouterModule } from '@angular/router';
import { HomeComponent } from './components/smart/home/home.component';
import { AppRoutingModule } from './app-routing.module';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { highlightProvider } from './modules/shared/highlight';
import { MonacoEditorModule } from '@materia-ui/ngx-monaco-editor';
import { ComponentReuseStrategy } from './modules/shared/component-reuse.strategy';
import { windowProvider, WindowToken } from './window';

@Injectable()
export class UnstripTrailingSlashLocation extends Location {
  public static stripTrailingSlash(url: string): string {
    return url;
  }
}

@NgModule({
  declarations: [HomeComponent],
  imports: [
    CommonModule,
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    RouterModule,
    MonacoEditorModule,
    // routing loads last
    AppRoutingModule,
  ],
  providers: [
    {
      provide: Location,
      useClass: UnstripTrailingSlashLocation,
    },
    highlightProvider(),
    { provide: RouteReuseStrategy, useClass: ComponentReuseStrategy },
    { provide: WindowToken, useFactory: windowProvider },
  ],
  bootstrap: [HomeComponent],
})
export class AppModule {}
