/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

import { HIGHLIGHT_OPTIONS } from 'ngx-highlightjs';

const languages = () => {
  return {
    json: () => import('highlight.js/lib/languages/json'),
    yaml: () => import('highlight.js/lib/languages/yaml'),
  };
};

export const highlightProvider = () => ({
  provide: HIGHLIGHT_OPTIONS,
  useValue: {
    languages: languages(),
  },
});
