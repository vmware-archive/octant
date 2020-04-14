import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ButtonGroupDemoComponent } from './button-group.demo';
import { ApiButtonGroupDemoComponent } from './api-button-group.demo';
import { AngularButtonGroupDemoComponent } from './angular-button-group.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { MarkdownService, MarkedOptions } from 'ngx-markdown';
import { ClarityModule } from '@clr/angular';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    ClarityModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: ButtonGroupDemoComponent }]),
  ],
  declarations: [
    AngularButtonGroupDemoComponent,
    ButtonGroupDemoComponent,
    ApiButtonGroupDemoComponent,
  ],
  providers: [MarkdownService, MarkedOptions],
  exports: [ButtonGroupDemoComponent],
})
export class ButtonGroupDemoModule {}
