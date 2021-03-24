import { Component, EventEmitter, Output } from '@angular/core';
import { ActionField, StepItem, StepperView } from '../../../models/content';
import {
  FormBuilder,
  FormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { ActionService } from '../../../services/action/action.service';

interface Choice {
  label: string;
  value: string;
  checked: boolean;
}

@Component({
  selector: 'app-stepper',
  templateUrl: './stepper.component.html',
  styleUrls: ['./stepper.component.scss'],
})
export class StepperComponent extends AbstractViewComponent<StepperView> {
  @Output()
  submit: EventEmitter<FormGroup> = new EventEmitter(true);

  @Output()
  cancel: EventEmitter<boolean> = new EventEmitter(true);

  formGroup: FormGroup;
  action: string;
  steps: StepItem[] = [];

  needParams = {
    min: true,
    max: true,
    minLength: true,
    maxLength: true,
    pattern: true,
    required: false,
    requiredTrue: false,
    email: false,
    nullValidator: false,
  };

  constructor(
    private formBuilder: FormBuilder,
    private actionService: ActionService
  ) {
    super();
  }

  update() {
    this.action = this.v.config.action;
    this.steps = this.v.config.steps;

    const stepGroups = this.createStepGroups(this.steps);
    this.formGroup = this.formBuilder.group(stepGroups);
  }

  createStepGroups(steps: StepItem[]): { [name: string]: any } {
    const stepGroups: { [name: string]: any } = {};

    steps?.forEach(step => {
      const controls = this.createControlGroups(step);
      stepGroups[step.name] = this.formBuilder.group(controls);
    });

    return stepGroups;
  }

  createControlGroups(step: StepItem): { [name: string]: any } {
    const controls: { [name: string]: any } = {};

    step.form?.fields?.forEach(field => {
      controls[field.name] = [
        field.value,
        this.getValidators(field.validators),
      ];
    });

    return controls;
  }

  onFormSubmit() {
    this.actionService.perform({
      action: this.action,
      ...this.formGroup.value,
    });
  }

  onFormCancel() {
    this.cancel.emit(true);
  }

  trackByFn(index, _) {
    return index;
  }

  getValidators(validators: { string: any }): ValidatorFn[] {
    if (!validators) {
      return [];
    }

    const vFn: ValidatorFn[] = [];
    const keys = Object.keys(validators);
    for (const key of keys) {
      const value = validators[key];

      // Check if function is expected
      if (this.needParams[key] === undefined) {
        console.error(`Unknown validation function ${key} for form`);
        continue;
      }

      // Verify how many params needs
      if (this.needParams[key]) {
        vFn.push(Validators[key](value));
      } else {
        vFn.push(Validators[key]);
      }
    }

    return vFn;
  }

  fieldChoices(field: ActionField) {
    return field.configuration.choices as Choice[];
  }
}
