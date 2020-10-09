/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

import { InjectionToken } from '@angular/core';

/**
 * WindowToken is an injection token that allows for injecting a `window` instance.
 */
export const WindowToken = new InjectionToken('WindowToken');

/**
 * windowProvider is the default window provider. It provides window from the browser.
 */
export function windowProvider(): any {
  return window;
}
