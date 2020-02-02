// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Location } from '@angular/common';
import { HttpClientModule } from '@angular/common/http';
import { Injectable, NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { AppRoutingModule } from './app-routing.module';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { highlightProvider } from './shared/highlight';

@Injectable()
export class UnstripTrailingSlashLocation extends Location {
  public static stripTrailingSlash(url: string): string {
    return url;
  }
}

@NgModule({
  declarations: [HomeComponent],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    RouterModule,

    // routing loads last
    AppRoutingModule,
  ],
  providers: [
    {
      provide: Location,
      useClass: UnstripTrailingSlashLocation,
    },
    highlightProvider(),
  ],
  bootstrap: [HomeComponent],
})
export class AppModule {}
