// Copyright (c) 2020 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { ClrDatagridComparatorInterface } from '@clr/angular';
import {
  TableRowWithMetadata,
  TimestampView,
} from '../modules/shared/models/content';

export class TimestampComparator
  implements ClrDatagridComparatorInterface<TableRowWithMetadata> {
  compare(a: TableRowWithMetadata, b: TableRowWithMetadata) {
    const rowA = a.data.Age as TimestampView;
    const rowB = b.data.Age as TimestampView;
    return rowA.config.timestamp - rowB.config.timestamp;
  }
}
