/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

import {
  AfterViewChecked,
  ChangeDetectionStrategy,
  Component,
  ViewEncapsulation,
} from '@angular/core';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import JSONEditor from 'jsoneditor';
import { JSONEditorView } from '../../../models/content';

@Component({
  selector: 'app-view-json',
  templateUrl: 'json-editor.component.html',
  styleUrls: ['./json-editor.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
  encapsulation: ViewEncapsulation.None,
})
export class JSONEditorComponent
  extends AbstractViewComponent<JSONEditorView>
  implements AfterViewChecked
{
  container: HTMLElement;
  id: string;
  options = {
    mode: '',
    mainMenuBar: false,
    navigationBar: false,
  };
  content: any;

  constructor() {
    super();
    this.id = Math.random().toString(36).substring(2, 15);
  }

  ngAfterViewChecked() {
    this.container = document.getElementById(this.id);
    if (this.container) {
      const editor = new JSONEditor(this.container, this.options);
      editor.set(this.content);

      if (this.v?.config?.collapsed) {
        editor.collapseAll();
      }
    }
  }

  update() {
    const view = this.v;
    this.options.mode = view.config.mode;

    this.isValidJson(view.config.content);
  }

  isValidJson(content: string): any {
    try {
      this.content = JSON.parse(content);
      return this.content;
    } catch (e) {
      this.options.mode = 'preview';
      this.content = {
        error: 'cannot parse json',
      };
      return this.content;
    }
  }
}
