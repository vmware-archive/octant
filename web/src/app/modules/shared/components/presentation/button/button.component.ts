import {
  Component,
  EventEmitter,
  Output,
  SecurityContext,
} from '@angular/core';
import {
  ButtonView,
  Confirmation,
  View,
  ModalView,
} from '../../../models/content';

import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { parse } from 'marked';
import { DomSanitizer } from '@angular/platform-browser';
import { ActionService } from '../../../services/action/action.service';
import { ModalService } from '../../../services/modal/modal.service';

@Component({
  selector: 'app-button',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.scss'],
})
export class ButtonComponent extends AbstractViewComponent<ButtonView> {
  @Output() buttonLoad: EventEmitter<boolean> = new EventEmitter(true);

  isModalOpen = false;
  modalTitle = '';
  modalBody = '';
  payload = {};
  class = '';

  modalView: View;

  constructor(
    private actionService: ActionService,
    private modalService: ModalService,
    private sanitize: DomSanitizer
  ) {
    super();
  }

  update() {
    if (this.v.config.confirmation) {
      this.class = 'btn btn-danger-outline btn-sm';
    } else {
      this.class = 'btn btn-outline btn-sm';
    }
    if (this.v.config.modal) {
      this.modalView = this.v.config.modal;
      const modal = this.modalView as ModalView;
      this.modalService.setState(modal.config.opened);
    }
  }

  onClick(payload: {}, confirmation?: Confirmation, modal?: View) {
    if (modal) {
      this.modalService.openModal();
    }
    if (confirmation) {
      this.activateModal(payload, confirmation);
    } else {
      this.buttonLoad.emit(true);
      this.doAction(payload);
    }
  }

  cancelModal() {
    this.resetModal();
  }

  acceptModal() {
    const payload = this.payload;
    this.resetModal();
    this.doAction(payload);
  }

  trackByFn(index, item) {
    return index;
  }

  private doAction(payload: {}) {
    this.actionService.perform(payload);
  }

  private activateModal(payload: {}, confirmation: Confirmation) {
    this.modalTitle = confirmation.title;
    this.modalBody = this.sanitize.sanitize(
      SecurityContext.HTML,
      parse(confirmation.body)
    );
    this.isModalOpen = true;

    this.payload = payload;
  }

  private resetModal() {
    this.isModalOpen = false;
    this.modalBody = '';
    this.modalTitle = '';
    this.payload = {};
  }
}
