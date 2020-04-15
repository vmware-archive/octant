// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//
export interface NavigationChild {
  module?: string;
  title: string;
  path: string;
  children?: NavigationChild[];
  iconName?: string;
  iconSource?: string;
  isLoading: boolean;
}

export interface Navigation {
  sections: NavigationChild[];
  defaultPath: string;
}
