// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { untilDestroyed } from 'ngx-take-until-destroy';
import {
  Filter,
  LabelFilterService,
} from 'src/app/modules/shared/services/label-filter/label-filter.service';

@Component({
  selector: 'app-filters',
  templateUrl: './filters.component.html',
  styleUrls: ['./filters.component.scss'],
})
export class FiltersComponent implements OnInit, OnDestroy {
  filters: Filter[];

  constructor(
    private labelFilter: LabelFilterService,
    private router: Router,
    private activatedRoute: ActivatedRoute
  ) {}

  ngOnInit() {
    this.labelFilter.filters.pipe(untilDestroyed(this)).subscribe(filters => {
      this.filters = filters;
      const filterParams = filters.map(filter =>
        encodeURIComponent(`${filter.key}:${filter.value}`)
      );
      const queryParams: Params = {
        filter: filterParams,
      };

      this.router.navigate([], {
        relativeTo: this.activatedRoute,
        queryParams,
        queryParamsHandling: 'merge',
      });
    });
  }

  ngOnDestroy() {}

  identifyFilter(index: number, item: Filter): string {
    return `${item.key}-${item.value}`;
  }

  remove(filter: Filter) {
    this.labelFilter.remove(filter);
  }
}
