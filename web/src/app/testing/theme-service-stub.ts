// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import {
  ThemeService,
  ThemeType,
} from '../modules/shared/services/theme/theme.service';
import { BehaviorSubject } from 'rxjs';

export const themeServiceStub: Partial<ThemeService> = {
  loadCSS: () => void 0,
  loadTheme: () => void 0,
  isLightThemeEnabled: () => true,
  switchTheme: () => void 0,
  themeType: new BehaviorSubject<ThemeType>('light'),
};
