import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { LogsDemoComponent } from './logs.demo';
import { ApiLogsDemoComponent } from './api-logs.demo';
import { AngularLogsDemoComponent } from './angular-logs.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: LogsDemoComponent }]),
  ],
  declarations: [
    AngularLogsDemoComponent,
    LogsDemoComponent,
    ApiLogsDemoComponent,
  ],
  exports: [LogsDemoComponent],
})
export class LogsDemoModule {}
