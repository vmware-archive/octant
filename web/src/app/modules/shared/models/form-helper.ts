import {
  FormArray,
  FormBuilder,
  FormControl,
  FormGroup,
  ValidatorFn,
  Validators,
} from '@angular/forms';
import { ActionField, ActionForm } from './content';

export interface Choice {
  label: string;
  value: string;
  checked: boolean;
}

// Check parameter necessary for validation functions
const validationNeedParams = {
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

// Class responsible to create a Form Group and add Validations Functions to form control
export class FormHelper {
  createFromGroup(form: ActionForm, formBuilder: FormBuilder): FormGroup {
    if (!form) {
      return;
    }

    const controls: { [name: string]: any } = {};
    form.fields.forEach(field => {
      if (field?.config.type === 'layout') {
        field?.config.configuration.fields.forEach(f => {
          this.createControls(controls, f);
        });
      } else {
        this.createControls(controls, field);
      }
    });
    return formBuilder.group(controls);
  }

  createControls(controls: { [name: string]: any }, field: ActionField) {
    controls[field.config.name] = [
      this.transformValue(field),
      this.getValidators(field.config.validators),
    ];

    if (
      field.config?.configuration?.choices &&
      (field.config.type === 'checkbox' || field.config.type === 'radio')
    ) {
      const choices: Choice[] = field.config.configuration.choices;
      controls[field.config.name] = new FormArray([]);
      choices.forEach((choice: Choice) => {
        if (choice.checked) {
          controls[field.config.name].push(new FormControl(choice.value));
        }
      });
    }
  }

  transformValue(field: ActionField): any {
    if (field.config.type === 'number') {
      if (field.config.value === '') {
        return null;
      }
      const value = +field.config.value;
      return Number.isNaN(value) ? 0 : value;
    }
    return field.config.value;
  }

  // Receive a hash with the validation name and the expected
  // params and return and array of functions
  getValidators(validators: { string: any }): ValidatorFn[] {
    if (!validators) {
      return [];
    }

    const vFn: ValidatorFn[] = [];
    const keys = Object.keys(validators);
    for (const key of keys) {
      const value = validators[key];

      // Check if function is expected
      if (validationNeedParams[key] === undefined) {
        console.error(`Unknown validation function ${key} for form`);
        continue;
      }

      // Verify how many params needs
      if (validationNeedParams[key]) {
        vFn.push(Validators[key](value));
      } else {
        vFn.push(Validators[key]);
      }
    }

    return vFn;
  }
}
