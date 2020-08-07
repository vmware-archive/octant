/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// GENERATED: do not edit!

import { ComponentFactory, FactoryMetadata } from './component-factory';
import { Component } from './component';

export interface LinkConfig {
  value: string;
  ref: string;
  status?: number;
  statusDetail?: Component<any>;
}

export interface LinkOptions {
  status?: number;
  statusDetail?: ComponentFactory<any>;
}

interface LinkParameters {
  value: string;
  ref: string;
  options?: LinkOptions;
  factoryMetadata?: FactoryMetadata;
}

export class LinkFactory implements ComponentFactory<LinkConfig> {
  private readonly value: string;
  private readonly ref: string;
  private readonly status: number | undefined;
  private readonly statusDetail: ComponentFactory<any> | undefined;
  private readonly factoryMetadata: FactoryMetadata | undefined;

  constructor({ value, ref, options, factoryMetadata }: LinkParameters) {
    this.value = value;
    this.ref = ref;
    this.factoryMetadata = factoryMetadata;

    if (options) {
      this.status = options.status;
      this.statusDetail = options.statusDetail;
    }
  }

  toComponent(): Component<LinkConfig> {
    return {
      metadata: {
        type: 'link',
        ...(this.factoryMetadata && { metadata: this.factoryMetadata }),
      },
      config: {
        value: this.value,
        ref: this.ref,

        ...(this.status && { status: this.status }),
        ...(this.statusDetail && {
          statusDetail: this.statusDetail.toComponent(),
        }),
      },
    };
  }
}
