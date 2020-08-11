import { Component } from '@angular/core';
import { CardListView, CardView } from '../../../models/content';
import { ViewService } from '../../../services/view/view.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-card-list',
  templateUrl: './card-list.component.html',
  styleUrls: ['./card-list.component.scss'],
})
export class CardListComponent extends AbstractViewComponent<CardListView> {
  constructor(private viewService: ViewService) {
    super();
  }

  update() {}

  identifyCard = (index: number, item: CardView): string => {
    return [index, this.viewService.viewTitleAsText(item)].join(',');
  };
}
