import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AnnotationsDemoComponent } from './annotations.demo';
import { ApiAnnotationsDemoComponent } from './api-annotations.demo';
import { AngularAnnotationsDemoComponent } from './angular-annotations.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { ClarityModule } from '@clr/angular';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    ClarityModule,
    SharedModule,
    CommonModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: AnnotationsDemoComponent }]),
  ],
  declarations: [
    AngularAnnotationsDemoComponent,
    AnnotationsDemoComponent,
    ApiAnnotationsDemoComponent,
  ],
  exports: [AnnotationsDemoComponent],
})
export class AnnotationsDemoModule {}
