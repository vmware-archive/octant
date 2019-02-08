import { compareModel } from '.'

describe('compareModel', () => {
  test('compare models that declare themselves as comparable', () => {
    const view1 = {
      type: 'text',
      value: 'a',
      isComparable: true,
    }

    expect(() => compareModel(view1, view1)).not.toThrow()
  })

  test('throw error when trying to compare incomparable', () => {
    const view1 = {
      type: 'text',
      value: 'a',
    }

    expect(() => compareModel(view1, view1)).toThrow()
  })

  test('throw error when trying to compare models with different types', () => {
    const view1 = {
      type: 'type1',
      value: 'a',
    }
    const view2 = {
      type: 'type2',
      value: 'a',
    }

    expect(() => compareModel(view1, view2)).toThrow()
  })
})
