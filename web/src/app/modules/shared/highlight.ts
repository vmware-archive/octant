/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import json from 'highlight.js/lib/languages/json';
import yaml from 'highlight.js/lib/languages/yaml';
import { HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';

export const hljsLanguages = () => [
  { name: 'yaml', func: yaml },
  { name: 'json', func: json },
];

export const highlightProvider = () => {
  return {
    provide: HIGHLIGHT_OPTIONS,
    useValue: {
      languages: hljsLanguages,
    },
  };
};
