import { Point } from './point';
import { Vector } from './vector';

describe('Point', () => {
  describe('toString', () => {
    it('converts a point to a string', () => {
      const point = new Point(1, 2);
      expect(point.toString()).toBe('1,2');
    });
  });

  describe('project', () => {
    it('projects a new point given a vector', () => {
      const point = new Point(0, 0);
      const vector: Vector = {
        magnitude: Math.sqrt(5 * 5 + 5 * 5),
        angle: 45,
      };

      const expected = new Point(5, 5);
      const got = point.project(vector);
      expect(got.x).toBeCloseTo(expected.x);
      expect(got.y).toBeCloseTo(expected.y);

    });
  });
});
