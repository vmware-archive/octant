// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnInit, Output } from '@angular/core';
import { ActionField, ActionForm } from '../../../models/content';
import {
  FormBuilder,
  FormGroup,
  ValidatorFn,
  Validators,
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

  formGroup: FormGroup;

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit() {
    if (this.form) {
      const controls: { [name: string]: any } = {};
      this.form.fields.forEach(field => {
        controls[field.name] = [
          field.value,
          this.getValidators(field.validators),
        ];
      });

      this.formGroup = this.formBuilder.group(controls);
    }
  }

  getValidators(validators: string[]): ValidatorFn[] {
    if (validators) {
      const vFn: ValidatorFn[] = [];
      validators.forEach(v => {
        vFn.push(Validators[v]);
      });
      return vFn;
    }
    return [];
  }

  fieldChoices(field: ActionField) {
    return field.configuration.choices as Choice[];
  }

  trackByFn(index, _) {
    return index;
  }
}
