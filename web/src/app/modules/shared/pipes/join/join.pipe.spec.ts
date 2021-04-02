import { JoinPipe } from './join.pipe';

describe('JoinPipe', () => {
  const pipe = new JoinPipe();

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  describe('empty values', () => {
    it('returns empty string if null', () => {
      expect(pipe.transform(null, 'some-delimiter')).toEqual('');
    });

    it('returns empty string if empty array', () => {
      expect(pipe.transform([], 'some-delimiter')).toEqual('');
    });
  });

  describe('single value', () => {
    it('returns only the single value', () => {
      expect(pipe.transform(['a-cool-value'], 'some-delimiter')).toEqual(
        'a-cool-value'
      );
    });
  });

  describe('multiple values', () => {
    it('returns joined values', () => {
      expect(pipe.transform(['valueA', 'valueB', 'valueC'], ',')).toEqual(
        'valueA,valueB,valueC'
      );
    });
  });
});
