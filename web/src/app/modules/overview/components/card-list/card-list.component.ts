import { Component, Input } from '@angular/core';
import { CardListView, CardView } from '../../../../models/content';
import { ViewService } from '../../services/view/view.service';

@Component({
  selector: 'app-view-card-list',
  templateUrl: './card-list.component.html',
  styleUrls: ['./card-list.component.scss'],
})
export class CardListComponent {
  @Input()
  view: CardListView;

  constructor(private viewService: ViewService) {}

  identifyCard = (index: number, item: CardView): string => {
    return [index, this.viewService.viewTitleAsText(item)].join(',');
  };
}
