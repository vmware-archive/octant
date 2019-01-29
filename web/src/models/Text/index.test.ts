import { compareTextModel } from '.'

describe('Text', () => {
  describe('compare', () => {
    const view1 = {
      type: 'text',
      value: 'a',
      title: '',
      isComparable: true,
    }
    const view2 = {
      type: 'text',
      value: 'b',
      title: '',
      isComparable: true,
    }

    test('less', () => {
      expect(compareTextModel(view1, view2)).toEqual(-1)
    })
    test('greater', () => {
      expect(compareTextModel(view2, view1)).toEqual(1)
    })
    test('equal', () => {
      expect(compareTextModel(view1, view1)).toEqual(0)
    })
  })
})
