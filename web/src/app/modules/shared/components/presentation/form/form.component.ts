// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, Input, OnInit } from '@angular/core';
import { ActionForm } from '../../../models/content';
import { FormBuilder, FormGroup } from '@angular/forms';
import { FormHelper } from '../../../models/form-helper';

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
    const formHelper = new FormHelper();
    this.formGroup = formHelper.createFromGroup(this.form, this.formBuilder);
  }
}
