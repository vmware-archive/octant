// Copyright (c) 2019 the Octant contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
//

import * as dagreD3 from 'dagre-d3';

export interface ResourceObject {
  name: string;
  apiVersion: string;
  kind: string;
  status: string;
}

export class ResourceNode {
  constructor(
    private readonly id: string,
    private readonly object: ResourceObject,
    private readonly isSelected: boolean
  ) {}

  toDescriptor(): dagreD3.Label {
    let nodeClass = `node-${this.object.status}`;
    if (this.isSelected) {
      nodeClass += ` selected`;
    }

    return {
      id: this.id,
      label: `${this.title()}${this.subTitle()}`,
      labelType: 'html',
      class: `${nodeClass}`,
    };
  }

  title(): string {
    return `<div class="resource-name">${this.object.name}</div>`;
  }

  subTitle(): string {
    return `<div class="resource-type">${this.object.apiVersion} ${this.object.kind}</div>`;
  }
}
