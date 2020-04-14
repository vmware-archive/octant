import { Component } from '@angular/core';
import {
  CardView,
  CardListView,
  TextView,
} from '../../../../../../../src/app/modules/shared/models/content';

const text: TextView = {
  config: {
    value: 'Card Title',
  },
  metadata: {
    type: 'text',
  },
};

const bodyView: TextView = {
  config: {
    value: 'body text',
  },
  metadata: {
    type: 'text',
  },
};

const view: CardView = {
  config: {
    body: bodyView,
    actions: null,
    alert: null,
  },
  metadata: {
    title: [text],
    type: 'card',
  },
};

const code = `card := component.NewCard([]component.TitleComponent{component.NewText("Card Title")})
card.SetBody(component.NewText("body text"))
`;

const json = JSON.stringify(view, null, 4);

@Component({
  selector: 'app-angular-card-demo',
  templateUrl: './angular-card.demo.html',
})
export class AngularCardDemoComponent {
  view = view;
  code = code;
  json = json;
}
