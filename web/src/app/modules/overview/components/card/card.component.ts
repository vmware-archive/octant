import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { Action, CardView, View } from '../../../../models/content';
import { FormGroup } from '@angular/forms';
import { ActionService } from '../../services/action/action.service';
import { ViewService } from '../../services/view/view.service';
import { FormComponent } from '../form/form.component';

@Component({
  selector: 'app-view-card',
  templateUrl: './card.component.html',
  styleUrls: ['./card.component.scss'],
})
export class CardComponent implements OnChanges {
  @Input()
  view: CardView;

  @ViewChild('appForm') appForm: FormComponent;

  title: string;

  body: View;

  currentAction: Action;

  constructor(
    private actionService: ActionService,
    private viewService: ViewService
  ) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.view.currentValue) {
      const view: CardView = changes.view.currentValue;
      if (view) {
        this.title = this.viewService.viewTitleAsText(view);
        this.body = view.config.body;
      }
    }
  }

  onActionSubmit(formGroup: FormGroup) {
    if (formGroup && formGroup.value) {
      this.actionService.perform(formGroup.value).subscribe();
      this.currentAction = undefined;
    }
  }

  onActionCancel() {
    this.currentAction = undefined;
  }

  setAction(action: Action) {
    this.currentAction = action;
  }

  trackByFn(index, item) {
    return index;
  }
}
