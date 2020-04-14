import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { SelectorsDemoComponent } from './selectors.demo';
import { ApiSelectorsDemoComponent } from './api-selectors.demo';
import { AngularSelectorsDemoComponent } from './angular-selectors.demo';

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
    RouterModule.forChild([{ path: '', component: SelectorsDemoComponent }]),
  ],
  declarations: [
    AngularSelectorsDemoComponent,
    SelectorsDemoComponent,
    ApiSelectorsDemoComponent,
  ],
  exports: [SelectorsDemoComponent],
  providers: [MarkdownService, MarkedOptions],
})
export class SelectorsDemoModule {}
