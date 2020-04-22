/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';

export function hljsLanguages() {
  return {
    json: () => import('highlight.js/lib/languages/json'),
    yaml: () => import('highlight.js/lib/languages/yaml'),
  };
}

export function highlightProvider() {
  return {
    provide: HIGHLIGHT_OPTIONS,
    useValue: {
      languages: hljsLanguages(),
    },
  };
}
