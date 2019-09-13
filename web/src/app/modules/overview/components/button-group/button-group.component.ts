import { Component, Input, OnInit } from '@angular/core';
import { ButtonGroupView, Confirmation } from '../../../../models/content';
import { ActionService } from '../../services/action/action.service';

@Component({
  selector: 'app-button-group',
  templateUrl: './button-group.component.html',
  styleUrls: ['./button-group.component.scss'],
})
export class ButtonGroupComponent implements OnInit {
  @Input() view: ButtonGroupView;

  isModalOpen = false;
  modalTitle = '';
  modalBody = '';
  payload = {};

  constructor(private actionService: ActionService) {}

  ngOnInit() {}

  onClick(payload: {}, confirmation?: Confirmation) {
    console.log({ payload, confirmation });
    if (confirmation) {
      this.activateModal(payload, confirmation);
    } else {
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
