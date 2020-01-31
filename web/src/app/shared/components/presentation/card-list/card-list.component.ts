import { Component, Input } from '@angular/core';
import { CardListView, CardView, View } from '../../../models/content';

@Component({
  selector: 'app-view-card-list',
  templateUrl: './card-list.component.html',
  styleUrls: ['./card-list.component.scss'],
})
export class CardListComponent {
  v: CardListView;

  @Input() set view(v: View) {
    this.v = v as CardListView;
  }
  get view() {
    return this.v;
  }

  constructor() {}

  identifyCard = (index: number, _: CardView): number => {
    return index;
  };
}
