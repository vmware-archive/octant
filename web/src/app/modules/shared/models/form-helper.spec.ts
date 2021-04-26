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
