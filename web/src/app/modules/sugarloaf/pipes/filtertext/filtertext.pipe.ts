// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { Pipe, PipeTransform } from '@angular/core';
import { Filter } from '../../../shared/services/label-filter/label-filter.service';

@Pipe({
  name: 'filtertext',
})
export class FilterTextPipe implements PipeTransform {
  transform(filters: Filter[]): string {
    if (filters && filters.length > 0) {
      return `Filter by labels (${filters.length} applied)`;
    }
    return 'Filter by labels';
  }
}
