// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { environment } from 'src/environments/environment';

// TODO: (bryanl) convert this to a service or remove it

export default function getAPIBase(): string {
  if (environment.production) {
    return window.location.origin;
  }
  return 'http://localhost:7777';
}
