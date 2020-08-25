import { Component, EventEmitter, Output } from '@angular/core';
import { ButtonGroupView, Confirmation } from '../../../models/content';
import { ActionService } from '../../../services/action/action.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';

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

  constructor(private actionService: ActionService) {
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
      });
    }
  }

  onClick(payload: {}, confirmation?: Confirmation) {
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
    this.modalBody = confirmation.body;
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
