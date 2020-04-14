import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CodeDemoComponent } from './code.demo';
import { ApiCodeDemoComponent } from './api-code.demo';
import { AngularCodeDemoComponent } from './angular-code.demo';
import { UtilsModule } from '../../../utils/utils.module';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { MarkdownService, MarkedOptions } from 'ngx-markdown';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: CodeDemoComponent }]),
  ],
  providers: [MarkdownService, MarkedOptions],
  declarations: [
    AngularCodeDemoComponent,
    CodeDemoComponent,
    ApiCodeDemoComponent,
  ],
  exports: [CodeDemoComponent],
})
export class CodeDemoModule {}
