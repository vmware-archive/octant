import { Component, Input, OnInit } from '@angular/core';
import { ActionField, ActionForm } from '../../models/content';
import { FormArray, FormControl, FormGroup } from '@angular/forms';
import { Choice } from '../../models/form-helper';
import trackByIndex from 'src/app/util/trackBy/trackByIndex';

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
}
