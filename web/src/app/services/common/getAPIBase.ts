// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { environment } from 'src/environments/environment';

// TODO: API_BASE this should be configurable
export default function getAPIBase(): string {
  if (environment.production) {
    return window.location.origin;
  }
  return 'http://localhost:7777';
}
