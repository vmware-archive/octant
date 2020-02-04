/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from './components/smart/home/home.component';
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '../shared/shared.module';
import { ClarityModule } from '@clr/angular';

const routes: Routes = [
  {
    path: '**',
    component: HomeComponent,
  },
];

@NgModule({
  imports: [
    CommonModule,
    ClarityModule,
    SharedModule,

    // routing must come last
    RouterModule.forChild(routes),
  ],
  declarations: [HomeComponent],
})
export class DenaliModule {}
