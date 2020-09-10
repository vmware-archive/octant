import { Component, ViewChild } from '@angular/core';
import { Action, CardView, TitleView, View } from '../../../models/content';
import { ActionService } from '../../../services/action/action.service';
import { FormComponent } from '../form/form.component';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

@Component({
  selector: 'app-view-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.scss'],
})
export class CardComponent extends AbstractViewComponent<CardView> {
  @ViewChild('appForm') appForm: FormComponent;

  title: TitleView[];

  body: View;

  currentAction: Action;

  constructor(private actionService: ActionService) {
    super();
  }

  update() {
    this.title = this.v.metadata.title as TitleView[];
    this.body = this.v.config.body;
  }

  onActionSubmit() {
    if (this.appForm?.formGroup && this.appForm?.formGroup.value) {
      this.actionService.perform(this.appForm.formGroup.value);
      this.currentAction = undefined;
    }
  }

  onActionCancel() {
    this.currentAction = undefined;
  }

  setAction(action: Action) {
    this.currentAction = action;
  }

  trackByFn(index, _) {
    return index;
  }
}
