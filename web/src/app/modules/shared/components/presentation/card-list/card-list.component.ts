import { Component, Input } from '@angular/core';
import { CardListView, CardView, View } from '../../../models/content';
import { ViewService } from '../../../services/view/view.service';

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

  constructor(private viewService: ViewService) {}

  identifyCard = (index: number, item: CardView): string => {
    return [index, this.viewService.viewTitleAsText(item)].join(',');
  };
}
