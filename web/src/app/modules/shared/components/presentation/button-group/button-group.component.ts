import {
  Component,
  EventEmitter,
  Output,
  SecurityContext,
} from '@angular/core';
import {
  ButtonGroupView,
  Confirmation,
  View,
  ModalView,
} from '../../../models/content';
import { ActionService } from '../../../services/action/action.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { ModalService } from '../../../services/modal/modal.service';
import { DomSanitizer } from '@angular/platform-browser';
import { parse } from 'marked';

@Component({
  selector: 'app-button-group',
  templateUrl: './button-group.component.html',
  styleUrls: ['./button-group.component.scss'],
})
export class ButtonGroupComponent extends AbstractViewComponent<
  ButtonGroupView
> {
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
    if (this.v.config.buttons) {
      this.v.config.buttons.forEach(button => {
        if (button.confirmation) {
          this.class = 'btn-danger-outline btn-sm';
        } else {
          this.class = 'btn-outline btn-sm';
        }
        if (button.modal) {
          this.modalView = button.modal;
          const modal = this.modalView as ModalView;
          this.modalService.setState(modal.config.opened);
        }
      });
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
