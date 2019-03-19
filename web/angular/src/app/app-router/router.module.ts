import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { ContentComponent } from '../layout/content/content.component';
import { PageNotFoundComponent } from '../util/page-not-found/page-not-found.component';

export const appRoutes: Routes = [
  { path: 'content', children: [{ path: '**', component: ContentComponent }] },
  { path: '', redirectTo: '/content/overview/namespace/default/', pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent },
];

@NgModule({
  declarations: [],
  imports: [RouterModule.forRoot(appRoutes), CommonModule],
})
export class AppRouterModule {}
