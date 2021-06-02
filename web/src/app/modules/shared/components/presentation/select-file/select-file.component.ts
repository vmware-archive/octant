// Copyright (c) 2021 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import {
  Component,
  ElementRef,
  EventEmitter,
  isDevMode,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { SelectFileView } from '../../../models/content';
import { WebsocketService } from '../../../../../data/services/websocket/websocket.service';
import { ActionService } from '../../../services/action/action.service';

@Component({
  selector: 'app-view-select-file',
  templateUrl: './select-file.component.html',
})
export class SelectFileComponent
  extends AbstractViewComponent<SelectFileView>
  implements OnInit {
  label: string;
  multiple: boolean;
  layout: string;
  status: string;
  statusMessage: string;
  action: string;

  @ViewChild('fileInput') fileInput: ElementRef;
  @Output() fileChanged: EventEmitter<any> = new EventEmitter<any>();

  constructor(private actionService: ActionService) {
    super();
  }

  update() {
    const view = this.v;
    this.label = view.config.label;
    this.layout = view.config.layout;
    this.multiple = view.config.multiple;
    this.status = view.config.status;
    this.statusMessage = view.config.statusMessage;
    this.action = view.config.action;
  }

  inputFileChanged(event) {
    if (event.target.files && event.target.files[0]) {
      if (this.fileChanged) {
        this.fileChanged.emit(event.target.files);
      }
      if (isDevMode()) {
        console.log('Selected file(s):', event.target.files);
      }

      if (this.action) {
        this.actionService.perform({
          action: this.action,
          files: event.target.files,
        });
      }
    }
  }

  reset() {
    this.fileInput.nativeElement.value = '';
    this.fileInput.nativeElement.dispatchEvent(new Event('change'));
  }
}
