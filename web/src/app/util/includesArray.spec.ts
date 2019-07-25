import { includesArray } from './includesArray';

describe('includesArray', () => {
  it('should return false if second param array is larger', () => {
    expect(includesArray(['a', 'b', 'c'], ['a', 'b', 'c', 'd'])).toBe(false);
  });

  it('should return true if second array is included in first array', () => {
    expect(includesArray(['a', 'b', 'c', 'd'], ['a', 'b', 'c'])).toBe(true);
  });

  it('should return false if second array is not included in first array', () => {
    expect(includesArray(['a', 'b', 'c', 'd'], ['a', 'e', 'c'])).toBe(false);
  });

  it('should return false if unordered second array is not included in first array', () => {
    expect(includesArray(['a', 'b', 'c', 'd'], ['a', 'c', 'b'])).toBe(false);
  });
});
