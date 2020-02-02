// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OverviewComponent } from './modules/sugarloaf/components/smart/overview/overview.component';

export const appRoutes: Routes = [
  {
    path: 'denali',
    loadChildren: () =>
      import('./modules/denali/denali.module').then(m => m.DenaliModule),
  },
  {
    path: '**',
    loadChildren: () =>
      import('./modules/sugarloaf/sugarloaf.module').then(
        m => m.SugarloafModule
      ),
  },
];

@NgModule({
  declarations: [],
  imports: [
    // routing must come last
    RouterModule.forRoot(appRoutes, {
      useHash: true,
      enableTracing: false,
    }),
  ],
})
export class AppRoutingModule {}
