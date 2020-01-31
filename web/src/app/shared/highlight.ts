/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import json from 'highlight.js/lib/languages/json';
import yaml from 'highlight.js/lib/languages/yaml';

export function hljsLanguages() {
  return [
    { name: 'yaml', func: yaml },
    { name: 'json', func: json },
  ];
}
