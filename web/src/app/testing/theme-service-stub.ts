// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
import { ThemeService } from '../modules/shared/services/theme/theme.service';

export const themeServiceStub: Partial<ThemeService> = {
  loadCSS: () => void 0,
  loadTheme: () => void 0,
  isLightThemeEnabled: () => true,
  onChange: () => void 0,
  offChange: () => void 0,
  switchTheme: () => void 0,
};
