// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClarityModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { OlympusStacksComponent } from './components/olympus-stacks/olympus-stacks.component';
import { WorkloadCardComponent } from './components/workload-card/workload-card.component';
import { WorkloadListComponent } from './components/workload-list/workload-list.component';

@NgModule({
  declarations: [
    OlympusStacksComponent,
    WorkloadCardComponent,
    WorkloadListComponent,
  ],
  imports: [
    CommonModule,
    ClarityModule,
    FormsModule,
  ]
})
export class OlympusStacksModule { }
