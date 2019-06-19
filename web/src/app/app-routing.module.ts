// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { OverviewComponent } from './modules/overview/overview.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';

export const appRoutes: Routes = [
  { path: 'content', children: [{ path: '**', component: OverviewComponent }] },
  // TODO: we shouldn't assume that the default namespace is valid
  { path: '', redirectTo: '/content/overview/namespace/default/', pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent },
];

@NgModule({
  declarations: [],
  imports: [RouterModule.forRoot(appRoutes, {
    useHash: true,
  }), CommonModule],
})
export class AppRoutingModule {}
