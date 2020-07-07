/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';
import { EditorView, View } from '../../../models/content';
import { ThemeService } from '../../../../sugarloaf/components/smart/theme-switch/theme-switch.service';
import { NamespaceService } from '../../../services/namespace/namespace.service';
import { ActionService } from '../../../services/action/action.service';

interface Options {
  readOnly: boolean;
  language: string;
  theme: string;
}

@Component({
  selector: 'app-view-editor',
  templateUrl: './editor.component.html',
  styleUrls: ['./editor.component.scss'],
})
export class EditorComponent implements OnChanges {
  private v: EditorView;

  @Input() set view(v: View) {
    this.v = v as EditorView;
  }

  get view() {
    return this.v;
  }

  set value(v: string) {
    if (v !== this.editorValue) {
      this.isModified = true;
    }
    this.editorValue = v;
  }

  get value() {
    return this.editorValue;
  }

  private editorValue: string;
  private pristineValue: string;
  private metadata: { [p: string]: string };

  isModified = false;

  options: Options;

  submitAction = 'action.octant.dev/update';
  submitLabel = 'Update';

  constructor(
    private themeService: ThemeService,
    private namespaceService: NamespaceService,
    private actionService: ActionService
  ) {
    this.options = {} as Options;
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes.view.currentValue) {
      const view = changes.view.currentValue as EditorView;

      if (!this.isModified) {
        this.editorValue = view.config.value;
        this.metadata = view.config.metadata;
        this.pristineValue = view.config.value;
        this.options.readOnly = view.config.readOnly;
        this.options.language = view.config.language;
      }

      this.submitAction = view.config.submitAction || this.submitAction;
      this.submitLabel = view.config.submitLabel || this.submitLabel;

      this.options.theme =
        this.themeService.currentType() === 'dark' ? 'vs-dark' : 'vs';
    }
  }

  submit() {
    const payload = {
      action: this.submitAction,
      update: this.value,
      ...(this.metadata || {
        namespace: this.namespaceService.activeNamespace.value,
      }),
    };
    this.actionService.perform(payload);
  }

  isUpdateEnabled() {
    return !this.isModified;
  }

  reset() {
    this.value = this.pristineValue;
  }
}
