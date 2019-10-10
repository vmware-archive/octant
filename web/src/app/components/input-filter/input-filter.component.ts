// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Component, HostListener, ElementRef, OnInit, OnDestroy } from '@angular/core';
import { untilDestroyed } from 'ngx-take-until-destroy';
import {
  Filter,
  LabelFilterService,
} from '../../services/label-filter/label-filter.service';

@Component({
  selector: 'app-input-filter',
  templateUrl: './input-filter.component.html',
  styleUrls: ['./input-filter.component.scss'],
})
export class InputFilterComponent implements OnInit, OnDestroy {
  inputValue = '';
  showTagList = false;
  filters: Filter[] = [];

  constructor(
    private eRef: ElementRef,
    private labelFilterService: LabelFilterService
  ) {}

  ngOnInit() {
    this.labelFilterService.filters
      .pipe(untilDestroyed(this))
      .subscribe(filters => {
        this.filters = filters;
      });
  }

  ngOnDestroy() {}

  @HostListener('document:click', ['$event'])
  outsideClick(event) {
    if (!this.eRef.nativeElement.contains(event.target)) {
      this.showTagList = false;
    }
  }

  toggleTagList() {
    this.showTagList = !this.showTagList;
  }

  identifyFilter(index: number, item: Filter): string {
    return `${item.key}-${item.value}`;
  }

  remove(filter: Filter) {
    this.labelFilterService.remove(filter);
  }

  get placeholderText(): string {
    if (this.filters && this.filters.length > 0) {
      return `Filter by labels (${this.filters.length} applied)`;
    } else {
      return 'Filter by labels';
    }
  }

  onEnter() {
    const filter = this.labelFilterService.decodeFilter(this.inputValue);
    if (filter) {
      this.labelFilterService.add(filter);
      this.inputValue = '';
      this.showTagList = true;
    } else {
      // TODO: user input value not a valid filter;
    }
  }

  clearAllFilters() {
    this.labelFilterService.clear();
    this.showTagList = false;
  }
}
