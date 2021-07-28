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
import { ActionService } from '../../../services/action/action.service';
import { ElectronService } from '../../../services/electron/electron.service';

type File = {
  name: string;
  type: string;
  path?: string;
  lastModified: number;
  size: number;
};

@Component({
  selector: 'app-view-select-file',
  templateUrl: './select-file.component.html',
})
export class SelectFileComponent
  extends AbstractViewComponent<SelectFileView>
  implements OnInit
{
  label: string;
  multiple: boolean;
  layout: string;
  status: string;
  statusMessage: string;
  action: string;

  @ViewChild('fileInput') fileInput: ElementRef;
  @Output() fileChanged: EventEmitter<any> = new EventEmitter<any>();

  constructor(
    private actionService: ActionService,
    private electronService: ElectronService
  ) {
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

  inputFileChanged(event: Event) {
    const files = (event.target as HTMLInputElement).files;
    const fileList: File[] = [];
    if (files && files[0]) {
      if (this.fileChanged) {
        this.fileChanged.emit(files);
      }
      if (isDevMode()) {
        console.log('Selected file(s):', files);
      }
      for (let i = 0; i < files.length; i++) {
        const file = files.item(i);
        let fileMetadata = {
          name: file.name,
          type: file.type,
          lastModified: file.lastModified,
          size: file.size,
        };

        if (this.electronService.isElectron()) {
          fileMetadata = { ...fileMetadata, ...{ path: file.path } };
        }
        fileList.push(fileMetadata);
      }

      if (this.action) {
        this.actionService.perform({
          action: this.action,
          files: fileList,
        });
      }
    }
  }

  reset() {
    this.fileInput.nativeElement.value = '';
    this.fileInput.nativeElement.dispatchEvent(new Event('change'));
  }
}
