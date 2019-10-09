// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OverviewComponent } from './modules/overview/overview.component';

export const appRoutes: Routes = [{ path: '**', component: OverviewComponent }];

@NgModule({
  declarations: [],
  imports: [
    RouterModule.forRoot(appRoutes, {
      useHash: true,
    }),
    CommonModule,
  ],
})
export class AppRoutingModule {}
