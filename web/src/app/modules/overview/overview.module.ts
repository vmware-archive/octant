// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { CommonModule } from '@angular/common';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ResizableModule } from 'angular-resizable-element';
import { OverviewComponent } from './overview.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SharedModule } from '../../shared/shared.module';

@NgModule({
  declarations: [OverviewComponent],
  imports: [BrowserAnimationsModule, CommonModule, SharedModule],
})
export class OverviewModule {}
