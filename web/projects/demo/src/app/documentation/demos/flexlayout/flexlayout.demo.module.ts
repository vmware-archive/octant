import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { FlexLayoutDemoComponent } from './flexlayout.demo';
import { ApiFlexLayoutDemoComponent } from './api-flexlayout.demo';
import { AngularFlexLayoutDemoComponent } from './angular-flexlayout.demo';
import { UtilsModule } from '../../../utils/utils.module';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: FlexLayoutDemoComponent }]),
  ],
  declarations: [
    AngularFlexLayoutDemoComponent,
    FlexLayoutDemoComponent,
    ApiFlexLayoutDemoComponent,
  ],
  exports: [FlexLayoutDemoComponent],
})
export class FlexLayoutDemoModule {}
