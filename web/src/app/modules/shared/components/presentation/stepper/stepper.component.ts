import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { ActionField, StepItem, StepperView } from '../../../models/content';
import { FormBuilder, FormGroup } from '@angular/forms';
import { AbstractViewComponent } from '../../abstract-view/abstract-view.component';
import { ActionService } from '../../../services/action/action.service';
import { FormHelper } from '../../../models/form-helper';

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
export class StepperComponent
  extends AbstractViewComponent<StepperView>
  implements OnInit
{
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

  ngOnInit() {
    this.formGroup = this.createStepGroups(this.steps);
  }

  // Each Step is a form, so it creates a form group for every form and encapsulate
  // them with another form group.
  createStepGroups(steps: StepItem[]): FormGroup {
    const stepGroups: { [name: string]: any } = {};
    const formHelper = new FormHelper();

    steps?.forEach(step => {
      stepGroups[step.name] = formHelper.createFromGroup(
        step.form,
        this.formBuilder
      );
    });

    return this.formBuilder.group(stepGroups);
  }

  update() {
    this.action = this.v.config.action;
    this.steps = this.v.config.steps;
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

  fieldChoices(field: ActionField) {
    return field.config.configuration.choices as Choice[];
  }

  formGroupFromName(step: StepItem): FormGroup {
    return this.formGroup.get(step.name) as FormGroup;
  }
}
