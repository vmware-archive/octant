import { Vector } from './vector';

export class Point {
  constructor(public x: number, public y: number)  {}

  toString() {
    return `${this.x},${this.y}`;
  }

  project(vector: Vector) {
    const radian = (Math.PI / 180) * vector.angle;
    const newX = this.x + Math.cos(radian) * vector.magnitude;
    const newY = this.y + Math.sin(radian) * vector.magnitude;

    return new Point(newX, newY);
  }
}
