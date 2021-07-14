import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { FormViewContainerComponent } from './form-view-container.component';
import { Component } from '@angular/core';
import {
  FormArray,
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { ActionForm } from '../../models/content';
import { FormHelper } from '../../models/form-helper';
import { CdsModule } from '@cds/angular';
import '@cds/core/radio/register.js';

@Component({
  template:
    '<app-form-view-container [form]="form" [formGroupContainer]="formGroup"></app-form-view-container>',
})
class TestWrapperComponent {
  form: ActionForm;
  formGroup: FormGroup;
}

describe('FormViewContainerComponent', () => {
  let component: TestWrapperComponent;
  let fixture: ComponentFixture<TestWrapperComponent>;
  let element: HTMLDivElement;
  let formHelper;

  const formBuilder: FormBuilder = new FormBuilder();

  beforeEach(
    waitForAsync(() => {
      TestBed.configureTestingModule({
        declarations: [TestWrapperComponent, FormViewContainerComponent],
        imports: [CdsModule, ReactiveFormsModule, FormsModule],
        providers: [{ provide: FormBuilder, useValue: formBuilder }],
      }).compileComponents();
    })
  );

  beforeEach(() => {
    fixture = TestBed.createComponent(TestWrapperComponent);
    component = fixture.componentInstance;
    element = fixture.nativeElement;
    formHelper = new FormHelper();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('form group', () => {
    it('creates radio', () => {
      component.form = {
        fields: [
          {
            config: {
              configuration: {
                choices: [
                  { label: 'a', value: 'a', checked: true },
                  { label: 'b', value: 'b', checked: false },
                  { label: 'c', value: 'c', checked: false },
                ],
              },
              label: 'label',
              name: 'name',
              type: 'radio',
              value: null,
              placeholder: '',
              error: null,
              validators: null,
            },
            metadata: { type: 'formField' },
          },
        ],
      };
      component.formGroup = formHelper.createFromGroup(
        component.form,
        formBuilder
      );
      fixture.detectChanges();
      expect(element.querySelector('cds-radio-group')).not.toBeNull();
      expect(element.querySelector('cds-radio')).not.toBeNull();
    });

    it('should create a select and verify is selected', () => {
      const name = 'name';
      component.form = {
        fields: [
          {
            config: {
              configuration: {
                choices: [
                  { label: 'a', value: 'a', checked: true },
                  { label: 'b', value: 'b', checked: false },
                  { label: 'c', value: 'c', checked: false },
                ],
              },
              label: 'label',
              name,
              type: 'select',
              value: null,
              placeholder: '',
              error: null,
              validators: null,
            },
            metadata: { type: 'formField' },
          },
        ],
      };
      component.formGroup = formHelper.createFromGroup(
        component.form,
        formBuilder
      );
      fixture.detectChanges();

      const selected = (component.formGroup.get(
        name
      ) as FormArray).getRawValue();
      expect(selected[0]).toEqual('a');
      expect(element.querySelector('cds-select')).not.toBeNull();
    });

    it('should create a select and verify is NOT selected', () => {
      const name = 'name';
      component.form = {
        fields: [
          {
            config: {
              configuration: {
                choices: [
                  { label: 'd', value: 'd', checked: false },
                  { label: 'b', value: 'b', checked: false },
                  { label: 'c', value: 'c', checked: false },
                ],
              },
              label: 'label',
              name,
              type: 'select',
              value: null,
              placeholder: '',
              error: null,
              validators: null,
            },
            metadata: { type: 'formField' },
          },
        ],
      };
      component.formGroup = formHelper.createFromGroup(
        component.form,
        formBuilder
      );
      fixture.detectChanges();

      const selected = (component.formGroup.get(
        name
      ) as FormArray).getRawValue();
      expect(selected[0]).toEqual(undefined);
      expect(element.querySelector('cds-select')).not.toBeNull();
    });
  });
});
