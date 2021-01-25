// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

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

    it('projects a new point given a vector and radian angle', () => {
      const point = new Point(10, 10);
      const vector: Vector = {
        magnitude: (10 * 2) / Math.sqrt(2),
        angle: (Math.PI / 180) * 45,
      };

      const expected = new Point(20, 20);
      const got = point.projectRadian(vector);
      expect(got.x).toBeCloseTo(expected.x);
      expect(got.y).toBeCloseTo(expected.y);
    });
  });
});
