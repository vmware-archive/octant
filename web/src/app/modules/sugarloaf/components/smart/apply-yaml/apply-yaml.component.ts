/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import { Component, HostListener, OnInit } from '@angular/core';
import { EditorView } from 'src/app/modules/shared/models/content';

@Component({
  selector: 'app-apply-yaml',
  templateUrl: './apply-yaml.component.html',
  styleUrls: ['./apply-yaml.component.scss'],
})
export class ApplyYAMLComponent implements OnInit {
  isOpen: boolean;
  editorView: EditorView = {
    config: {
      value: '',
      language: 'yaml',
      readOnly: false,
      metadata: null,
      submitAction: 'action.octant.dev/apply',
      submitLabel: 'Apply',
    },
    metadata: {
      type: 'editor',
      title: [],
    },
  };

  constructor() {}

  ngOnInit() {
    this.isOpen = false;
  }

  @HostListener('window:keydown', ['$event'])
  keyEvent(event: KeyboardEvent) {
    if (event.ctrlKey && event.key === 'y') {
      event.preventDefault();
      event.cancelBubble = true;
      this.isOpen = !this.isOpen;
    }
  }

  toggleModal() {
    this.isOpen = true;
  }
}
