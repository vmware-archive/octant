import { compareTimestampModel } from './timestamp'

describe('Timestamp', () => {
  describe('compare', () => {
    const ts1 = {
      type: 'timestamp',
      isComparable: true,
      timestamp: 5,
    }
    const ts2 = {
      type: 'timestamp',
      isComparable: true,
      timestamp: 7,
    }

    test('less', () => {
      expect(compareTimestampModel(ts1, ts2)).toEqual(-1)
    })
    test('greater', () => {
      expect(compareTimestampModel(ts2, ts1)).toEqual(1)
    })
    test('equal', () => {
      expect(compareTimestampModel(ts1, ts1)).toEqual(0)
    })
  })
})
