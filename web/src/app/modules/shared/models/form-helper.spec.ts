import { FormHelper } from './form-helper';

describe('FormHelper', () => {
  let formHelper: FormHelper;

  beforeEach(() => {
    formHelper = new FormHelper();
  });

  it('converts number', () => {
    expect(
      formHelper.transformValue({
        metadata: {
          type: 'form',
        },
        config: {
          configuration: {},
          value: 3,
          name: '',
          type: 'number',
          error: null,
          label: null,
          placeholder: null,
          validators: null,
          width: 0,
        },
      })
    ).toEqual(3);
  });

  it('converts stringed number', () => {
    expect(
      formHelper.transformValue({
        metadata: {
          type: 'form',
        },
        config: {
          configuration: null,
          value: '123',
          name: '',
          type: 'number',
          error: null,
          label: null,
          placeholder: null,
          validators: null,
          width: 0,
        },
      })
    ).toEqual(123);
  });

  it('converts text', () => {
    expect(
      formHelper.transformValue({
        metadata: {
          type: 'form',
        },
        config: {
          configuration: null,
          value: 'hello',
          name: '',
          type: 'text',
          error: null,
          label: null,
          placeholder: null,
          validators: null,
        },
      })
    ).toEqual('hello');
  });

  it('converts NaN', () => {
    expect(
      formHelper.transformValue({
        metadata: {
          type: 'form',
        },
        config: {
          configuration: null,
          value: NaN,
          name: '',
          type: 'number',
          error: null,
          label: null,
          placeholder: null,
          validators: null,
        },
      })
    ).toEqual(0);
  });
});

describe('FormHelper:createControls', () => {
  let formHelper: FormHelper;

  beforeEach(() => {
    formHelper = new FormHelper();
  });

  it('should create a form array', () => {
    const control = {};
    const name = 'test';
    formHelper.createControls(control, {
      metadata: {
        type: 'form',
      },
      config: {
        configuration: {
          choices: [
            { label: 'a', value: 'a', checked: true },
            { label: 'b', value: 'b', checked: false },
            { label: 'c', value: 'c', checked: false },
          ],
        },
        value: 3,
        name,
        type: 'select',
        error: null,
        label: null,
        placeholder: null,
        validators: null,
        width: 0,
      },
    });

    expect(control[name]).not.toBeNull();
    expect(control[name].pristine).toBeTruthy();
  });

  it('should create an array', () => {
    const control = {};
    const name = 'test';
    formHelper.createControls(control, {
      metadata: {
        type: 'form',
      },
      config: {
        configuration: {
          choices: [
            { label: 'a', value: 'a', checked: true },
            { label: 'b', value: 'b', checked: false },
            { label: 'c', value: 'c', checked: false },
          ],
        },
        value: 3,
        name,
        type: 'number',
        error: null,
        label: null,
        placeholder: null,
        validators: null,
        width: 0,
      },
    });

    expect(control[name]).not.toBeNull();
    expect(control[name].pristine).toBe(undefined);
  });
});
