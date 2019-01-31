import { compareModel } from '.'

describe('compareModel', () => {
  test('compare models that declare themselves as comparable', () => {
    const view1 = {
      type: 'text',
      value: 'a',
      title: '',
      isComparable: true,
    }

    expect(() => compareModel(view1, view1)).not.toThrow()
  })

  test('throw error when trying to compare uncomparable', () => {
    const view1 = {
      type: 'text',
      value: 'a',
      title: '',
    }

    expect(() => compareModel(view1, view1)).toThrow()
  })

  test('throw error when trying to compare models with different types', () => {
    const view1 = {
      type: 'type1',
      value: 'a',
      title: '',
    }
    const view2 = {
      type: 'type2',
      value: 'a',
      title: '',
    }

    expect(() => compareModel(view1, view2)).toThrow()
  })
})
