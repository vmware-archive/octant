// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'default',
  pure: true,
})
export class DefaultPipe implements PipeTransform {
  transform(value: any, defaultValue: any): any {
    return value || defaultValue;
  }
}
