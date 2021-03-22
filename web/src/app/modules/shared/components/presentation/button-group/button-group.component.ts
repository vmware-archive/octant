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
  Button,
} from '../../../models/content';
import { ActionService } from '../../../services/action/action.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { ModalService } from '../../../services/modal/modal.service';
import { DomSanitizer } from '@angular/platform-browser';
import { parse } from 'marked';

// Todo: move to common import with link.component.ts
const isUrlExternal = url =>
  url?.indexOf('://') > 0 || url?.indexOf('//') === 0;

@Component({
  selector: 'app-button-group',
  templateUrl: './button-group.component.html',
  styleUrls: ['./button-group.component.scss'],
})
export class ButtonGroupComponent extends AbstractViewComponent<ButtonGroupView> {
  @Output() buttonLoad: EventEmitter<boolean> = new EventEmitter(true);

  isModalOpen = false;
  modalTitle = '';
  modalBody = '';
  payload = {};
  buttons: Button[];

  modalView: View;

  constructor(
    private actionService: ActionService,
    private modalService: ModalService,
    private sanitize: DomSanitizer
  ) {
    super();
  }

  update() {
    const current = this.v;
    this.buttons = current.config.buttons;
    if (this.buttons) {
      this.buttons.forEach(button => {
        if (button.modal) {
          this.modalView = button.modal;
          const modal = this.modalView as ModalView;
          this.modalService.setState(modal.config.opened);
        }
        if (button.confirmation) {
          button.style = 'btn-danger-outline btn-sm';
        } else {
          button.style = this.buttonClass(
            button.style,
            button.size,
            button.status
          );
        }
      });
    }
  }

  onClick(
    payload: {},
    confirmation?: Confirmation,
    modal?: View,
    ref?: string
  ) {
    if (modal) {
      this.modalService.openModal();
    }
    if (confirmation) {
      this.activateModal(payload, confirmation);
    } else {
      this.buttonLoad.emit(true);
      this.doAction(payload);
    }
    if (ref) {
      if (isUrlExternal(ref)) {
        window.open(ref);
      } else {
        location.href = ref;
      }
    }
  }

  buttonClass(style: string, size: string, status: string): string {
    const buttonBase = 'btn';
    const result = [buttonBase];
    // Default size and style is btn-outline btn-sm
    let buttonSize = 'btn';
    let buttonStyles = 'btn';

    if (size) {
      buttonSize += '-' + size;
    }
    if (size === 'lg') {
      buttonSize = '';
      result.push(buttonSize);
    } else if (buttonSize.length > 3) {
      // no-op
      result.push(buttonSize);
    } else {
      result.push(buttonSize + '-' + 'sm');
    }

    if (status) {
      if (status !== 'disabled') {
        buttonStyles += '-' + status;
      }
      if (status === 'disabled') {
        result.push('disabled');
      }
    }
    if (style) {
      buttonStyles += '-' + style;
    }
    if (buttonStyles.length > 3) {
      result.push(buttonStyles);
    } else {
      result.push(buttonStyles + '-' + 'outline');
    }
    return result.filter(String).join(' ');
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
