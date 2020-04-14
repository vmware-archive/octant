import { Component } from '@angular/core';
import {
  CardListView,
  CardView,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const text1: TextView = {
  config: {
    value: 'Card 1',
  },
  metadata: {
    type: 'text',
  },
};

const text2: TextView = {
  config: {
    value: 'Card 2',
  },
  metadata: {
    type: 'text',
  },
};

const bodyView: TextView = {
  config: {
    value: 'card content',
  },
  metadata: {
    type: 'text',
  },
};

const card1: CardView = {
  config: {
    body: bodyView,
    actions: null,
    alert: null,
  },
  metadata: {
    title: [text1],
    type: 'text',
  },
};

const card2: CardView = {
  config: {
    body: bodyView,
    actions: null,
    alert: null,
  },
  metadata: {
    title: [text2],
    type: 'text',
  },
};

const cardListView: CardListView = {
  config: {
    cards: [card1, card2],
  },
  metadata: {
    type: 'cards',
  },
};

const code = `card list component
`;

@Component({
  selector: 'app-angular-card-list-demo',
  templateUrl: './angular-card-list.demo.html',
})
export class AngularCardListDemoComponent {
  view = cardListView;
  code = code;
}
