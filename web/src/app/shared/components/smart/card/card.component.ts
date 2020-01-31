import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { Action, CardView, TitleView, View } from '../../../models/content';
import { FormGroup } from '@angular/forms';
import { ActionService } from '../../../../modules/overview/services/action/action.service';
import { FormComponent } from '../../presentation/form/form.component';

@Component({
  selector: 'app-view-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.scss'],
})
export class CardComponent implements OnChanges {
  private v: CardView;

  @Input() set view(v: View) {
    this.v = v as CardView;
  }

  get view() {
    return this.v;
  }

  @ViewChild('appForm') appForm: FormComponent;

  title: TitleView[];

  body: View;

  currentAction: Action;

  constructor(private actionService: ActionService) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view: CardView = changes.view.currentValue;
      if (view) {
        this.title = view.metadata.title as TitleView[];
        this.body = view.config.body;
      }
    }
  }

  onActionSubmit(formGroup: FormGroup) {
    if (formGroup && formGroup.value) {
      this.actionService.perform(formGroup.value);
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
