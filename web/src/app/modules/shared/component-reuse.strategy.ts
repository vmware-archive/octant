/*
 *  Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import {
  ActivatedRouteSnapshot,
  DetachedRouteHandle,
  RouteReuseStrategy,
} from '@angular/router';
import { Injectable } from '@angular/core';

const genKey = (r: ActivatedRouteSnapshot) => {
  return `${r.url.join('/')}`;
};

@Injectable()
export class ComponentReuseStrategy implements RouteReuseStrategy {
  store(route: ActivatedRouteSnapshot, handle: DetachedRouteHandle): void {}

  retrieve(route: ActivatedRouteSnapshot): DetachedRouteHandle {
    return null;
  }
  shouldAttach(route: ActivatedRouteSnapshot): boolean {
    return false;
  }

  shouldDetach(route: ActivatedRouteSnapshot): boolean {
    return false;
  }

  /**
   * if navigating between tabs, reuse the route
   *
   * @param future future route
   * @param curr current route
   */
  shouldReuseRoute(
    future: ActivatedRouteSnapshot,
    curr: ActivatedRouteSnapshot
  ): boolean {
    return genKey(future) === genKey(curr);
  }
}
