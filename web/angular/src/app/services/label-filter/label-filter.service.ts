import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

export interface Filter {
  key: string;
  value: string;
}

/**
 * hash a string into a number
 *
 * @param s string
 */
const hashCode = (s: string): number => {
  let h: number;
  for (let i = (h = 0); i < s.length; i++) {
    // tslint:disable-next-line:no-bitwise
    h = (Math.imul(31, h) + s.charCodeAt(i)) | 0;
  }
  return h;
};


@Injectable({
  providedIn: 'root',
})
export class LabelFilterService {
  private labels: { [key: number]: Filter } = {};

  private filterObservable = new BehaviorSubject<Filter[]>([]);

  constructor() {}

  select(key: string, value: string) {
    const label = `${key}:${value}`;
    const hashed = hashCode(label);
    if (!this.labels.hasOwnProperty(hashed)) {
      this.labels[hashed] = {
        key,
        value,
      };
      this.publish();
    }
  }

  set(filters: Filter[]) {
    filters.forEach((filter) => {
      this.select(filter.key, filter.value);
    });
  }

  remove(filter: Filter) {
    const hash = this.hash(filter);
    delete this.labels[hash];
    this.publish();
  }

  filters(): Observable<Filter[]> {
    return this.filterObservable;
  }

  private publish() {
    this.filterObservable.next(Object.values(this.labels));
  }

  private hash(filter: Filter): number {
    const s = `${filter.key}:${filter.value}`;
    return hashCode(s);
  }
}
