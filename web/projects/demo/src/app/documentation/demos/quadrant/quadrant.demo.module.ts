import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { QuadrantDemoComponent } from './quadrant.demo';
import { ApiQuadrantDemoComponent } from './api-quadrant.demo';
import { AngularQuadrantDemoComponent } from './angular-quadrant.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { UtilsModule } from '../../../utils/utils.module';
import { ClarityModule } from '@clr/angular';
import { MarkdownService, MarkedOptions } from 'ngx-markdown';

@NgModule({
  imports: [
    UtilsModule,
    ClarityModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: QuadrantDemoComponent }]),
  ],
  declarations: [
    AngularQuadrantDemoComponent,
    QuadrantDemoComponent,
    ApiQuadrantDemoComponent,
  ],
  exports: [QuadrantDemoComponent],
  providers: [MarkdownService, MarkedOptions],
})
export class QuadrantDemoModule {}
