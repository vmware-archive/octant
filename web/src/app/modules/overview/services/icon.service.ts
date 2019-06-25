// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Injectable } from '@angular/core';

export interface IconAble {
  iconName?: string;
  iconSource?: string;
}

@Injectable({
  providedIn: 'root',
})
export class IconService {
  constructor() {}

  load(item: IconAble): string {
    if (!item.iconName || item.iconName === '') {
      return '';
    }

    // tslint:disable:no-string-literal
    const clarityIcons = window['ClarityIcons'];

    if (!clarityIcons.has(item.iconName)) {
      clarityIcons.add({ [item.iconName]: item.iconSource });
    }

    return item.iconName;
  }
}
