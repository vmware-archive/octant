/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import json from 'highlight.js/lib/languages/json';
import yaml from 'highlight.js/lib/languages/yaml';
import { HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';

export function hljsLanguages() {
  return {
    json: () => json,
    yaml: () => yaml,
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
