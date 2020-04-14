import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { PodStatusDemoComponent } from './pod-status.demo';
import { ApiPodStatusDemoComponent } from './api-pod-status.demo';
import { AngularPodStatusDemoComponent } from './angular-pod-status.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: PodStatusDemoComponent }]),
  ],
  declarations: [
    AngularPodStatusDemoComponent,
    PodStatusDemoComponent,
    ApiPodStatusDemoComponent,
  ],
  exports: [PodStatusDemoComponent],
})
export class PodStatusDemoModule {}
