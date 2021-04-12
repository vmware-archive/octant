import { Component, EventEmitter, Output } from '@angular/core';
import { ButtonGroupView } from '../../../models/content';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-button-group',
  templateUrl: './button-group.component.html',
  styleUrls: ['./button-group.component.scss'],
})
export class ButtonGroupComponent extends AbstractViewComponent<ButtonGroupView> {
  @Output() buttonLoad: EventEmitter<boolean> = new EventEmitter(true);

  constructor() {
    super();
  }

  update() {}

  trackByFn(index) {
    return index;
  }
}
