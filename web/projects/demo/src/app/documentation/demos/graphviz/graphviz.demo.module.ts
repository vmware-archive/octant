import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { GraphvizDemoComponent } from './graphviz.demo';
import { ApiGraphvizDemoComponent } from './api-graphviz.demo';
import { AngularGraphvizDemoComponent } from './angular-graphviz.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: GraphvizDemoComponent }]),
  ],
  declarations: [
    AngularGraphvizDemoComponent,
    GraphvizDemoComponent,
    ApiGraphvizDemoComponent,
  ],
  exports: [GraphvizDemoComponent],
})
export class GraphvizDemoModule {}
