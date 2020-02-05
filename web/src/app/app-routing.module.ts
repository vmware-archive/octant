// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

export const appRoutes: Routes = [
  {
    path: 'denali',
    loadChildren: './modules/denali/denali.module#DenaliModule',
  },
  {
    path: '**',
    loadChildren: './modules/sugarloaf/sugarloaf.module#SugarloafModule',
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
