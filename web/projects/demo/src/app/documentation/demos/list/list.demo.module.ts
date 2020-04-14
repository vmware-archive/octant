import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ListDemoComponent } from './list.demo';
import { ApiListDemoComponent } from './api-list.demo';
import { AngularListDemoComponent } from './angular-list.demo';
import { MarkdownService, MarkedOptions } from 'ngx-markdown';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: ListDemoComponent }]),
  ],
  providers: [MarkdownService, MarkedOptions],
  declarations: [
    AngularListDemoComponent,
    ListDemoComponent,
    ApiListDemoComponent,
  ],
  exports: [ListDemoComponent],
})
export class ListDemoModule {}
