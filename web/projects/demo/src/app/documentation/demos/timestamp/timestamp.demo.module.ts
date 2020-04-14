import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { TimestampDemoComponent } from './timestamp.demo';
import { ApiTimestampDemoComponent } from './api-timestamp.demo';
import { AngularTimestampDemoComponent } from './angular-timestamp.demo';

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
    RouterModule.forChild([{ path: '', component: TimestampDemoComponent }]),
  ],
  declarations: [
    AngularTimestampDemoComponent,
    TimestampDemoComponent,
    ApiTimestampDemoComponent,
  ],
  exports: [TimestampDemoComponent],
  providers: [MarkdownService, MarkedOptions],
})
export class TimestampDemoModule {}
