// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import { TextView, View } from '../../models/content';

export const titleAsText = (titleViews: View[]): string => {
  if (!titleViews) {
    return '';
  }
  // assume it's a text title for now
  return titleViews.map((titleView: TextView) => titleView.config.value).join(' / ');
};

export class ViewUtil {
  constructor(private view: View) {}

  titleAsText(): string {
    return titleAsText(this.view.metadata.title);
  }
}
