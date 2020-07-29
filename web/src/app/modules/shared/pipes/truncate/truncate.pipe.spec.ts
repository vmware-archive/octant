import { TruncatePipe } from './truncate.pipe';

describe('TruncatePipe', () => {
  const pipe = new TruncatePipe();

  it('create an instance', () => {
    expect(pipe).toBeTruthy();
  });

  it('more than 28 chars', () => {
    expect(pipe.transform('abcdefghijklmnopqrstuvwxyz123456789')).toEqual(
      'abcdefghij...vwxyz123456789'
    );
  });

  it('equal to 28 chars', () => {
    expect(pipe.transform('abcdefghijklmnopqrstuvwxyz12')).toEqual(
      'abcdefghijklmnopqrstuvwxyz12'
    );
  });

  it('less than 28 chars', () => {
    expect(pipe.transform('abc')).toEqual('abc');
  });
});
