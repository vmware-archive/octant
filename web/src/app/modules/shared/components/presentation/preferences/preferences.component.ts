import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  Output,
  SimpleChanges,
} from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup } from '@angular/forms';
import {
  Condition,
  Operation,
  PreferenceElement,
  Preferences,
} from '../../../models/preference';
import trackByIdentity from 'src/app/util/trackBy/trackByIdentity';
import { startWith } from 'rxjs/operators';

interface StringDict {
  [key: string]: string;
}

/**
 * Checks a condition against the current state.
 *
 * @param condition a condition
 * @param currentState the current state as a StringDict
 */
const checkCondition = (
  condition: Condition,
  currentState: StringDict
): boolean => {
  switch (condition.op) {
    case Operation.Equal:
      return currentState[condition.lhs] !== condition.rhs;
    default:
      // fail open
      return true;
  }
};

/**
 * Returns a list of elements defined in preferences.
 *
 * @param preferences preferences
 */
const elements = (preferences: Preferences): PreferenceElement[] => {
  return preferences.panels.reduce<PreferenceElement[]>((accum, panel) => {
    panel.sections.forEach(section => {
      accum.push(...section.elements);
    });
    return accum;
  }, []);
};

@Component({
  selector: 'app-preferences',
  templateUrl: './preferences.component.html',
  styleUrls: ['./preferences.component.scss'],
})
export class PreferencesComponent implements OnChanges {
  isOpenValue: boolean;

  @Input()
  get isOpen() {
    return this.isOpenValue;
  }

  set isOpen(v: boolean) {
    this.isOpenValue = v;
    this.isOpenChange.emit(this.isOpenValue);
  }

  @Input()
  preferences: Preferences;

  @Output()
  preferencesChanged = new EventEmitter();

  @Output()
  isOpenChange = new EventEmitter<boolean>();

  @Output()
  reset = new EventEmitter<void>();

  form: FormGroup = new FormGroup({});

  controls: { [key: string]: AbstractControl } = {};
  conditions: { [key: string]: Condition[] } = {};

  activeTabs: { [key: string]: boolean } = {};

  trackByIdentity = trackByIdentity;

  constructor(private fb: FormBuilder) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes.preferences && !changes.preferences.currentValue) {
      return;
    }

    const preferences = this.preferences;

    preferences.panels.forEach(
      (panel, i) => (this.activeTabs[panel.name] = i === 0)
    );

    elements(preferences).forEach(element => {
      this.controls[element.name] = this.fb.control({
        value: element.value,
        disabled: false,
      });
      if (element.disableConditions) {
        this.conditions[element.name] = element.disableConditions;
      }
    });

    this.form = this.fb.group(this.controls);

    this.form.valueChanges
      .pipe(startWith(this.form.getRawValue()))
      .subscribe((update: StringDict) => {
        this.onValueChanged(update);
      });
  }

  onCancel() {
    this.isOpen = false;
  }

  onDropDownValueChange(event, name) {
    this.form.value[name] = event;
  }

  onSubmit(): void {
    if (this.form.valid) {
      this.preferencesChanged.emit(this.form.value);
      this.isOpen = false;
    }
  }

  onReset(): void {
    this.reset.emit();
    this.isOpen = false;
  }

  private onValueChanged(update: StringDict) {
    Object.keys(this.controls).forEach(key => {
      if (!this.conditions[key]) {
        return;
      }
      if (this.checkConditions(this.conditions[key] as Condition[], update)) {
        if (this.controls[key].enabled) {
          this.controls[key].disable({ emitEvent: false });
        }
      } else {
        if (this.controls[key].disabled) {
          this.controls[key].enable({ emitEvent: false });
        }
      }
    });
  }

  private checkConditions(
    conditions: Condition[],
    current: StringDict
  ): boolean {
    return conditions
      .map<boolean>(condition => checkCondition(condition, current))
      .includes(true);
  }
}
