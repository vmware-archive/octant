import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import _ from 'lodash';
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router';

export interface Filter {
  key: string;
  value: string;
}

@Injectable({
  providedIn: 'root',
})
export class LabelFilterService {
  public filters = new BehaviorSubject<Filter[]>([]);
  private activatedRoute: ActivatedRoute;

  constructor(private router: Router) {
    this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.activatedRoute = this.router.routerState.root;
        const params = _.map(this.filters.getValue(), this.encodeFilter);
        this.router.navigate([], {
          relativeTo: this.activatedRoute,
          replaceUrl: true,
          queryParams: { filter: params },
          queryParamsHandling: 'merge',
        });
      }
    });

    this.router.routerState.root.queryParamMap.subscribe((paramMap) => {
      if (_.includes(paramMap.keys, 'filter')) {
        const filtersRaw = paramMap.getAll('filter');
        const filters = _.map(filtersRaw, this.decodeFilter);
        this.filters.next(filters);
      }
    });
  }

  add(filter: Filter): void {
    const current = this.filters.getValue();
    if (_.find(current, filter)) {
      return;
    }
    current.push(filter);
    this.publish(current);
  }

  remove(filter: Filter): void {
    const current = this.filters.getValue();
    _.remove(current, (f) => _.isEqual(filter, f));
    this.publish(current);
  }

  clearAll(): void {
    this.publish([]);
  }

  private encodeFilter(fil: Filter): string {
    return `${fil.key}:${fil.value}`;
  }

  public decodeFilter(filStr: string): Filter | null {
    const spl = filStr.split(':');
    if (spl.length > 1) {
      return { key: spl[0], value: spl[1] };
    }
    return null;
  }

  private publish(list: Filter[]): void {
    this.filters.next(list);
    const filterParams = list.map(this.encodeFilter);
    this.router.navigate([], {
      relativeTo: this.activatedRoute,
      queryParams: { filter: filterParams },
      queryParamsHandling: 'merge',
    });
  }
}
