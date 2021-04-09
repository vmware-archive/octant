import {
  FormArray,
  FormControl,
  ValidatorFn,
  Validators,
} from '@angular/forms';

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
  createFromGroup(form, formBuilder) {
    if (!form) {
      return;
    }

    const controls: { [name: string]: any } = {};
    form.fields.forEach(field => {
      controls[field.name] = [
        this.transformValue(field),
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

    return formBuilder.group(controls);
  }

  transformValue(field): any {
    if (field.type === 'number') {
      if (field.value === '') {
        return null;
      }
      const value = +field.value;
      return Number.isNaN(value) ? 0 : value;
    }
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
