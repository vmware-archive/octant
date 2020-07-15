/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package configuration

import "github.com/vmware-tanzu/octant/internal/describer"

var (
	pluginDescriber = NewPluginListDescriber()

	applyYamlDescriber = NewApplyYamlDescriber()

	rootDescriber = describer.NewSection(
		"/",
		"Configuration",
		pluginDescriber,
		applyYamlDescriber,
	)
)
