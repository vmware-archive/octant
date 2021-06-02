/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cluster

// InfoInterface provides connection details for a cluster
type InfoInterface interface {
	Context() string
	Cluster() string
	Server() string
	User() string
}
