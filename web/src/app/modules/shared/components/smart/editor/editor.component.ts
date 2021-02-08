/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import { Component, OnDestroy, OnInit } from '@angular/core';
import { EditorView } from '../../../models/content';
import { NamespaceService } from '../../../services/namespace/namespace.service';
import { ActionService } from '../../../services/action/action.service';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { ThemeService } from '../../../services/theme/theme.service';
import { Subscription } from 'rxjs';

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
export class EditorComponent
  extends AbstractViewComponent<EditorView>
  implements OnInit, OnDestroy {
  set value(v: string) {
    if (v !== this.editorValue) {
      this.isModified = true;
    }
    this.editorValue = v;
  }

  get value() {
    return this.editorValue;
  }

  private subscriptionTheme: Subscription;
  private syncMonacoTheme: () => void;
  private editorValue: string;
  private pristineValue: string;
  uri: string;
  private metadata: { [p: string]: string };

  isModified = false;

  options: Options = { theme: 'vs-dark', language: 'yaml', readOnly: false };

  submitAction = 'action.octant.dev/update';
  submitLabel = 'Update';

  constructor(
    private namespaceService: NamespaceService,
    private themeService: ThemeService,
    private actionService: ActionService
  ) {
    super();

    this.uri =
      'file:text-' + Math.random().toString(36).substring(2, 15) + '.yaml';

    this.syncMonacoTheme = () => {
      const theme = this.themeService.isLightThemeEnabled() ? 'vs' : 'vs-dark';
      this.options = { ...this.options, theme };
    };

    this.syncMonacoTheme();
  }

  ngOnInit() {
    this.subscriptionTheme = this.themeService.themeType.subscribe(() =>
      this.syncMonacoTheme()
    );
  }

  update() {
    const view = this.v;

    if (!this.isModified) {
      this.editorValue = view.config.value;
      this.metadata = view.config.metadata;
      this.pristineValue = view.config.value;
      this.options.readOnly = view.config.readOnly;
      this.options.language = view.config.language;
    }

    this.submitAction = view.config.submitAction || this.submitAction;
    this.submitLabel = view.config.submitLabel || this.submitLabel;
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

  ngOnDestroy() {
    this.subscriptionTheme.unsubscribe();
  }
}
