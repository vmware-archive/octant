/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package javascript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dop251/goja"
)

type httpClient struct {
	vm   *goja.Runtime
	this *goja.Object
}

// CreateHTTPClientObject creates an object that wraps HTTP client calls and exposes
// them as methods to be used in the JavaScript runtime.
func CreateHTTPClientObject(vm *goja.Runtime, this *goja.Object) goja.Value {
	client := vm.NewObject()
	h := &httpClient{
		vm:   vm,
		this: this,
	}
	if err := client.Set("get", h.get); err != nil {
		return vm.NewTypeError(fmt.Errorf("httpClient.Set.get: %w", err))
	}
	if err := client.Set("getJSON", h.getJSON); err != nil {
		return vm.NewTypeError(fmt.Errorf("httpClient.Set.getJSON: %w", err))
	}
	if err := client.Set("post", h.post); err != nil {
		return vm.NewTypeError(fmt.Errorf("httpClient.Set.post: %w", err))
	}
	return client
}

func (h *httpClient) httpGet(c goja.FunctionCall) (goja.Callable, []byte, error) {
	if len(c.Arguments) != 2 {
		return nil, nil, fmt.Errorf("invalid arguments")
	}

	urlArg := c.Argument(0).String()
	if urlArg == "" {
		return nil, nil, fmt.Errorf("empty url")
	}

	callbackArg := c.Argument(1).ToObject(h.vm)
	callback, ok := goja.AssertFunction(callbackArg)
	if !ok || callback == nil {
		return nil, nil, fmt.Errorf("bad callback function")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(urlArg)
	if err != nil {
		return nil, nil, fmt.Errorf("get: %w", err)
	}
	defer func() {
		_ = r.Body.Close()
	}()
	response, err := ioutil.ReadAll(r.Body)
	return callback, response, err
}

func (h *httpClient) get(c goja.FunctionCall) goja.Value {
	callback, response, err := h.httpGet(c)
	if err != nil {
		return h.vm.NewTypeError(fmt.Errorf("get: %w", err))
	}
	cr, err := callback(h.this, h.vm.ToValue(response))
	if err != nil {
		return h.vm.NewTypeError(fmt.Errorf("get: %w", err))
	}
	return cr
}

func (h *httpClient) getJSON(c goja.FunctionCall) goja.Value {
	callback, response, err := h.httpGet(c)
	if err != nil {
		return h.vm.NewTypeError(fmt.Errorf("getJSON: %w", err))
	}

	var target interface{}
	if err := json.NewDecoder(bytes.NewReader(response)).Decode(&target); err != nil {
		return h.vm.NewTypeError(fmt.Errorf("decoding: %w", err))
	}

	cr, err := callback(h.this, h.vm.ToValue(target))
	if err != nil {
		return h.vm.NewTypeError(fmt.Errorf("getJSON: %w", err))
	}
	return cr
}

func (h *httpClient) post(_ goja.FunctionCall) goja.Value {
	return h.vm.NewGoError(fmt.Errorf("not implemented"))
}
