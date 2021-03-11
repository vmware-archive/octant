/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package json

import jsoniter "github.com/json-iterator/go"

// Globally configure jsoniter to use fastest configuration together with sorting
var jsonConfig = jsoniter.Config{
	EscapeHTML:                    false,
	MarshalFloatWith6Digits:       true,
	ObjectFieldMustBeSimpleString: true,
	SortMapKeys:                   true,
}.Froze()

var NewEncoder = jsonConfig.NewEncoder
var NewDecoder = jsonConfig.NewDecoder
var Unmarshal = jsonConfig.Unmarshal
var Marshal = jsonConfig.Marshal
var MarshalIndent = jsonConfig.MarshalIndent
