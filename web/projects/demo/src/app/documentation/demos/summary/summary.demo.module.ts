import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { SummaryDemoComponent } from './summary.demo';
import { ApiSummaryDemoComponent } from './api-summary.demo';
import { AngularSummaryDemoComponent } from './angular-summary.demo';

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
    RouterModule.forChild([{ path: '', component: SummaryDemoComponent }]),
  ],
  declarations: [
    AngularSummaryDemoComponent,
    SummaryDemoComponent,
    ApiSummaryDemoComponent,
  ],
  exports: [SummaryDemoComponent],
  providers: [MarkdownService, MarkedOptions],
})
export class SummaryDemoModule {}
