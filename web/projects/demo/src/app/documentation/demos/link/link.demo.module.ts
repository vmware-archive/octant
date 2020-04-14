import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { LinkDemoComponent } from './link.demo';
import { ApiLinkDemoComponent } from './api-link.demo';
import { AngularLinkDemoComponent } from './angular-link.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';
import { UtilsModule } from '../../../utils/utils.module';

@NgModule({
  imports: [
    UtilsModule,
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: LinkDemoComponent }]),
  ],
  declarations: [
    AngularLinkDemoComponent,
    LinkDemoComponent,
    ApiLinkDemoComponent,
  ],
  exports: [LinkDemoComponent],
})
export class LinkDemoModule {}
