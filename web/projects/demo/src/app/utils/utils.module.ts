import { NgModule } from '@angular/core';
import { ClarityModule } from '@clr/angular';
import { CodeTabsComponent } from './code-tab.component';
import { SharedModule } from '../../../../../src/app/modules/shared/shared.module';
@NgModule({
  imports: [ClarityModule, SharedModule],
  declarations: [CodeTabsComponent],
  exports: [CodeTabsComponent],
})
export class UtilsModule {}
