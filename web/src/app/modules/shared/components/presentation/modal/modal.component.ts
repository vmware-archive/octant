import { Component, ViewChild, OnInit, OnDestroy } from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import {
  ActionField,
  TitleView,
  ModalView,
  View,
  ActionForm,
} from '../../../models/content';
import { ModalService } from '../../../services/modal/modal.service';
import { Subscription } from 'rxjs';
import {
  FormBuilder,
  FormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { ClrForm } from '@clr/angular';
import { WebsocketService } from '../../../services/websocket/websocket.service';

interface Choice {
  label: string;
  value: string;
  checked: boolean;
}

@Component({
  selector: 'app-view-modal',
  templateUrl: './modal.component.html',
  styleUrls: ['./modal.component.scss'],
})
export class ModalComponent
  extends AbstractViewComponent<ModalView>
  implements OnInit, OnDestroy {
  @ViewChild(ClrForm) clrForm: ClrForm;

  title: TitleView[];
  body: View;
  form: ActionForm;
  opened = false;
  size: string;
  formGroup: FormGroup;
  action: string;

  private modalSubscription: Subscription;

  constructor(
    private formBuilder: FormBuilder,
    private modalService: ModalService,
    private websocketService: WebsocketService
  ) {
    super();
  }

  ngOnInit() {
    this.modalSubscription = this.modalService.isOpened.subscribe(isOpened => {
      this.opened = isOpened;
    });
  }

  ngOnDestroy(): void {
    if (this.modalSubscription) {
      this.modalSubscription.unsubscribe();
    }
  }

  update() {
    this.title = this.v.metadata.title as TitleView[];
    this.body = this.v.config.body;
    this.size = this.v.config.size;
    this.form = this.v.config.form;
    this.opened = this.v.config.opened;
    this.modalService.setState(this.opened);

    if (this.form) {
      const controls: { [name: string]: any } = {};
      this.form.fields.forEach(field => {
        controls[field.name] = [
          field.value,
          this.getValidators(field.validators),
        ];
      });
      this.action = this.form.action;
      this.formGroup = this.formBuilder.group(controls);
    }
  }

  getValidators(validators: string[]): ValidatorFn[] {
    if (validators) {
      const vFn: ValidatorFn[] = [];
      validators.forEach(v => {
        vFn.push(Validators[v]);
      });
      return vFn;
    }
    return [];
  }

  trackByFn(index, _) {
    return index;
  }

  fieldChoices(field: ActionField) {
    return field.configuration.choices as Choice[];
  }

  onFormSubmit() {
    if (this.formGroup.invalid) {
      this.clrForm.markAsTouched();
    } else {
      this.websocketService.sendMessage('action.octant.dev/performAction', {
        action: this.action,
        formGroup: this.formGroup.value,
      });
    }
  }
}
