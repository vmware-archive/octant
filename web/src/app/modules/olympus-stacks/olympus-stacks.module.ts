// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClarityModule } from '@clr/angular';
import { FormsModule } from '@angular/forms';
import { OlympusStacksComponent } from './components/olympus-stacks/olympus-stacks.component';

@NgModule({
  declarations: [OlympusStacksComponent],
  imports: [
    CommonModule,
    ClarityModule,
    FormsModule,
  ]
})
export class OlympusStacksModule { }
