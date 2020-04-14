import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { LabelSelectorDemoComponent } from './label-selector.demo';
import { ApiLabelSelectorDemoComponent } from './api-label-selector.demo';
import { AngularLabelSelectorDemoComponent } from './angular-label-selector.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([
      { path: '', component: LabelSelectorDemoComponent },
    ]),
  ],
  declarations: [
    AngularLabelSelectorDemoComponent,
    LabelSelectorDemoComponent,
    ApiLabelSelectorDemoComponent,
  ],
  exports: [LabelSelectorDemoComponent],
})
export class LabelSelectorDemoModule {}
