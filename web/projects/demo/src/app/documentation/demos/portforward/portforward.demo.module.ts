import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { PortForwardDemoComponent } from './portforward.demo';
import { ApiPortForwardDemoComponent } from './api-portforward.demo';
import { AngularPortForwardDemoComponent } from './angular-portforward.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: PortForwardDemoComponent }]),
  ],
  declarations: [
    AngularPortForwardDemoComponent,
    PortForwardDemoComponent,
    ApiPortForwardDemoComponent,
  ],
  exports: [PortForwardDemoComponent],
})
export class PortForwardDemoModule {}
