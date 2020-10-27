import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import {
  ActionForm,
  Button,
  ModalView,
  TitleView,
  View,
} from '../../../models/content';
import { FormComponent } from '../form/form.component';
import { ModalService } from '../../../services/modal/modal.service';
import { Subscription } from 'rxjs';
import { ActionService } from '../../../services/action/action.service';

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
  @ViewChild('modalAppForm') modalAppForm: FormComponent;

  title: TitleView[];
  body: View;
  form: ActionForm;
  opened = false;
  size: string;
  action: string;
  buttons: Button[];

  private modalSubscription: Subscription;

  constructor(
    private actionService: ActionService,
    private modalService: ModalService
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
    this.action = this.form?.action;
    this.buttons = this.v.config.buttons;
  }

  onFormSubmit() {
    if (this.modalAppForm && this.modalAppForm.formGroup.valid) {
      this.actionService.perform({
        action: this.action,
        ...this.modalAppForm.formGroup.value,
      });
      this.opened = false;
    }
  }

  onClick(payload: {}) {
    this.actionService.perform(payload);
    this.opened = false;
  }
}
