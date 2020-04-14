import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ContainersDemoComponent } from './containers.demo';
import { ApiContainersDemoComponent } from './api-containers.demo';
import { AngularContainersDemoComponent } from './angular-containers.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: ContainersDemoComponent }]),
  ],
  declarations: [
    AngularContainersDemoComponent,
    ContainersDemoComponent,
    ApiContainersDemoComponent,
  ],
  exports: [ContainersDemoComponent],
})
export class ContainersDemoModule {}
