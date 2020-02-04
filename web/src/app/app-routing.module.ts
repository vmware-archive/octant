// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

const denaliModule = () =>
  import('./modules/denali/denali.module').then(m => m.DenaliModule);

const sugarloafModule = () =>
  import('./modules/sugarloaf/sugarloaf.module').then(m => m.SugarloafModule);

export const appRoutes: Routes = [
  {
    path: 'denali',
    loadChildren: denaliModule,
  },
  {
    path: '**',
    loadChildren: sugarloafModule,
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
