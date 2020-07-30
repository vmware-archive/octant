import {
  Component,
  OnInit,
  Input,
  SimpleChanges,
  EventEmitter,
  Output,
  OnChanges,
  OnDestroy,
  isDevMode,
} from '@angular/core';
import {
  StepperView,
  View,
  ActionField,
  StepItem,
} from '../../../models/content';
import {
  FormGroup,
  FormBuilder,
  Validators,
  ValidatorFn,
} from '@angular/forms';
import { WebsocketService } from '../../../services/websocket/websocket.service';

interface Choice {
  label: string;
  value: string;
  checked: boolean;
}

@Component({
  selector: 'app-stepper',
  templateUrl: './stepper.component.html',
  styleUrls: ['./stepper.component.sass'],
})
export class StepperComponent implements OnInit {
  private v: StepperView;

  @Input() set view(v: View) {
    this.v = v as StepperView;
  }

  get view() {
    return this.v;
  }

  @Output()
  submit: EventEmitter<FormGroup> = new EventEmitter(true);

  @Output()
  cancel: EventEmitter<boolean> = new EventEmitter(true);

  formGroup: FormGroup;
  action: string;
  steps: StepItem[];

  constructor(
    private formBuilder: FormBuilder,
    private websocketService: WebsocketService
  ) {
    this.steps = [] as StepItem[];
  }

  ngOnInit() {
    if (this.v) {
      this.action = this.v.config.action;
      this.steps = this.v.config.steps;

      const stepGroups = this.createStepGroups(this.steps);
      this.formGroup = this.formBuilder.group(stepGroups);
    }
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
    const form = JSON.stringify(this.formGroup.value);
    this.websocketService.sendMessage(this.action, form);
    if (isDevMode()) {
      console.log('stepper form: ' + form);
    }
  }

  onFormCancel() {
    this.cancel.emit(true);
  }

  trackByFn(index, item) {
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
