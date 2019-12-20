// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { environment } from 'src/environments/environment';

export default function getAPIBase(): string {
  if (environment.production) {
    return window.location.origin;
  }
  return 'http://localhost:7777';
}
