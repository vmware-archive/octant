// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnInit, Output } from '@angular/core';
import { ActionField, ActionForm } from '../../../models/content';
import {
  FormArray,
  FormBuilder,
  FormControl,
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
  formArray: FormArray;

  constructor(private formBuilder: FormBuilder) {}

  ngOnInit() {
    if (this.form) {
      const controls: { [name: string]: any } = {};
      this.form.fields.forEach(field => {
        controls[field.name] = [
          field.value,
          this.getValidators(field.validators),
        ];
        if (field.configuration?.choices && field.type === 'checkbox') {
          const choices: Choice[] = field.configuration.choices;
          controls[field.name] = new FormArray([]);
          choices.forEach((choice: Choice) => {
            if (choice.checked) {
              controls[field.name].push(new FormControl(choice.value));
            }
          });
        }
      });
      this.formGroup = this.formBuilder.group(controls);
    }
  }

  onCheck(event, field: string) {
    this.formArray = this.formGroup.get(field) as FormArray;
    if (event.target.checked) {
      this.formArray.push(new FormControl(event.target.value));
    } else {
      this.formArray.controls.forEach((fc: FormControl, index: number) => {
        if (fc.value === event.target.value) {
          this.formArray.removeAt(index);
        }
      });
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
