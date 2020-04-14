import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { TerminalDemoComponent } from './terminal.demo';
import { ApiTerminalDemoComponent } from './api-terminal.demo';
import { AngularTerminalDemoComponent } from './angular-terminal.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: TerminalDemoComponent }]),
  ],
  declarations: [
    AngularTerminalDemoComponent,
    TerminalDemoComponent,
    ApiTerminalDemoComponent,
  ],
  exports: [TerminalDemoComponent],
})
export class TerminalDemoModule {}
