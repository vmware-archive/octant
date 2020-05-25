/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package electron

import (
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

func platformWindowOptions(in astilectron.WindowOptions) astilectron.WindowOptions {
	in.TitleBarStyle = astikit.StrPtr("hiddenInset")

	return in
}
