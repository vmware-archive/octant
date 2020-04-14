import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { CardListDemoComponent } from './card-list.demo';
import { ApiCardListDemoComponent } from './api-card-list.demo';
import { AngularCardListDemoComponent } from './angular-card-list.demo';

import { SharedModule } from '../../../../../../../src/app/modules/shared/shared.module';

@NgModule({
  imports: [
    SharedModule,
    CommonModule,
    FormsModule,
    FormsModule,
    RouterModule.forChild([{ path: '', component: CardListDemoComponent }]),
  ],
  declarations: [
    AngularCardListDemoComponent,
    CardListDemoComponent,
    ApiCardListDemoComponent,
  ],
  exports: [CardListDemoComponent],
})
export class CardListDemoModule {}
