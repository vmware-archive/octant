/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

// GENERATED: do not edit!

/**
 * Metadata contains component metadata.
 */
export interface Metadata {
    /**
     * type is the type of component.
     */
    type: string;
    /*
     * title is the optional title of the component.
    */
    title?: Component<any>[];
    /*
     * accessor is the component's optional accessor.
    */
    accessor?: string;
}

/**
 * Component is a generic component.
 */
export interface Component<T> {
    /*
     * metadata is the component's metadata.
    */
    metadata: Metadata;
    /*
     * config is configuration for the component.
    */
    config: T;
}
