// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OverviewComponent } from './modules/overview/overview.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { NamespaceResolver } from 'src/app/services/namespace/namespace-resolver.service';

export const appRoutes: Routes = [
  { path: 'content', children: [{ path: '**', component: OverviewComponent }] },
  {
    path: '',
    component: OverviewComponent,
    resolve: {
      namespace: NamespaceResolver,
    },
    pathMatch: 'full',
  },
  { path: '**', component: PageNotFoundComponent },
];

@NgModule({
  declarations: [],
  imports: [
    RouterModule.forRoot(appRoutes, {
      useHash: true,
    }),
    CommonModule,
  ],
  providers: [NamespaceResolver],
})
export class AppRoutingModule {}
