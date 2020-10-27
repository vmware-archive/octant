import { Component, EventEmitter, isDevMode, Output } from '@angular/core';
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
}
