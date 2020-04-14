import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { IFrameDemoComponent } from './iframe.demo';
import { ApiIFrameDemoComponent } from './api-iframe.demo';
import { AngularIFrameDemoComponent } from './angular-iframe.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: IFrameDemoComponent }]),
  ],
  declarations: [
    AngularIFrameDemoComponent,
    IFrameDemoComponent,
    ApiIFrameDemoComponent,
  ],
  exports: [IFrameDemoComponent],
})
export class IFrameDemoModule {}
