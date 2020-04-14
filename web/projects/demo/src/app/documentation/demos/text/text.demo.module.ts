import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { TextDemoComponent } from './text.demo';
import { ApiTextDemoComponent } from './api-text.demo';
import { AngularTextDemoComponent } from './angular-text.demo';

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
    RouterModule.forChild([{ path: '', component: TextDemoComponent }]),
  ],
  declarations: [
    AngularTextDemoComponent,
    TextDemoComponent,
    ApiTextDemoComponent,
  ],
  exports: [TextDemoComponent],
  providers: [MarkdownService, MarkedOptions],
})
export class TextDemoModule {}
