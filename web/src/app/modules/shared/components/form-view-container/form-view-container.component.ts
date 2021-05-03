import { Component, Input, OnInit } from '@angular/core';
import { ActionField, ActionForm } from '../../models/content';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';

import '@cds/core/checkbox/register.js';
import '@cds/core/input/register.js';
import '@cds/core/textarea/register.js';
import '@cds/core/input/register.js';
import '@cds/core/radio/register.js';
import '@cds/core/select/register.js';
import { Choice } from '../../models/form-helper';

@Component({
  selector: 'app-form-view-container',
  templateUrl: './form-view-container.component.html',
  styleUrls: ['./form-view-container.component.scss'],
})
export class FormViewContainerComponent implements OnInit {
  @Input()
  form: ActionForm;
  @Input()
  formGroupContainer: FormGroup;

  formArray: FormArray;

  trackByFn = trackByIndex;

  ngOnInit(): void {}

  onCheck(event, field: string) {
    this.formArray = this.formGroupContainer.get(field) as FormArray;
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

  fieldChoices(field: ActionField) {
    return field.configuration.choices as Choice[];
  }

  isInvalid(fieldName: string) {
    const field = this.formGroupContainer.get(fieldName);
    return field.invalid && (field.dirty || field.touched);
  }
}
