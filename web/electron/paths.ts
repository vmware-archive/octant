/*
 *  Copyright (c) 2021 the Octant contributors. All Rights Reserved.
 *  SPDX-License-Identifier: Apache-2.0
 *
 */

import { app } from 'electron';
import * as path from 'path';
import * as os from 'os';

let date: string = new Date().toISOString();

export const tmpPath = path.join(os.tmpdir(), 'octant');
export const apiLogPath = path.join(tmpPath, 'api.out-' + date + '.log');
export const errLogPath = path.join(tmpPath, 'api.err-' + date + '.log');
export const iconPath = path.join(
  app.getAppPath(),
  'dist/octant/assets/icons/icon.png'
);
