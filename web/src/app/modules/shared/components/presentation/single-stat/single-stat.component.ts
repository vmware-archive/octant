import { Component } from '@angular/core';
import { SingleStatView } from '../../../models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-single-stat',
  templateUrl: './single-stat.component.html',
  styleUrls: ['./single-stat.component.scss'],
})
export class SingleStatComponent extends AbstractViewComponent<SingleStatView> {
  constructor() {
    super();
  }

  update() {}
}
