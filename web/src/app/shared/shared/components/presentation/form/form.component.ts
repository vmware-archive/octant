// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { ActionField, ActionForm } from '../../../../../models/content';
import {
  AbstractControl,
  FormBuilder,
  FormControl,
  FormGroup,
} from '@angular/forms';

interface Choice {
  label: string;
  value: string;
  checked: boolean;
}

@Component({
  selector: 'app-form',
  templateUrl: './form.component.html',
  styleUrls: ['./form.component.scss'],
})
export class FormComponent implements OnInit {
  @Input()
  form: ActionForm;

  @Input()
  title: string;

  @Output()
  submit: EventEmitter<FormGroup> = new EventEmitter(true);

  @Output()
  cancel: EventEmitter<boolean> = new EventEmitter(true);

  formGroup: FormGroup;

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit() {
    if (this.form) {
      const controls: { [name: string]: AbstractControl } = {};
      this.form.fields.forEach(field => {
        const value = field.value;
        controls[field.name] = new FormControl(value);
      });

      this.formGroup = this.formBuilder.group(controls);
    }
  }

  onFormSubmit() {
    this.submit.emit(this.formGroup);
  }

  onFormCancel() {
    this.cancel.emit(true);
  }

  fieldChoices(field: ActionField) {
    return field.configuration.choices as Choice[];
  }

  trackByFn(index, item) {
    return index;
  }
}
