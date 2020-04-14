import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { ExpressionSelectorDemoComponent } from './expression-selector.demo';
import { ApiExpressionSelectorDemoComponent } from './api-expression-selector.demo';
import { AngularExpressionSelectorDemoComponent } from './angular-expression-selector.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([
      { path: '', component: ExpressionSelectorDemoComponent },
    ]),
  ],
  declarations: [
    AngularExpressionSelectorDemoComponent,
    ExpressionSelectorDemoComponent,
    ApiExpressionSelectorDemoComponent,
  ],
  exports: [ExpressionSelectorDemoComponent],
})
export class ExpressionSelectorDemoModule {}
