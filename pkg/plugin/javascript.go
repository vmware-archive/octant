/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"sync"

	ocontext "github.com/vmware-tanzu/octant/internal/context"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/eventloop"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/plugin/javascript"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/log"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func IsJavaScriptPlugin(pluginName string) bool {
	return path.Ext(pluginName) == ".js"
}

// JSRuntimeFactory functions creates a JavaScript runtime for a JavaScript plugin.
type JSRuntimeFactory func(context.Context, string) (*eventloop.EventLoop, error)

// JSClassExtractor functions extract the default class from a runtime.
type JSClassExtractor func(*goja.Runtime) (*goja.Object, error)

// JSMetadataExtractor functions extract JavaScript plugin metadata from a runtime.
type JSMetadataExtractor func(*goja.Runtime, goja.Value) (*Metadata, error)

// WithRuntimeFactory option replaces the default JSRuntimeFactory function of a JSPlugin.
func WithRuntimeFactory(prf JSRuntimeFactory) func(*jsPlugin) {
	return func(js *jsPlugin) {
		js.runtimeFactory = prf
	}
}

// WithClassExtractor option replaces the default JSClassExtractor function of a JSPlugin.
func WithClassExtractor(pce JSClassExtractor) func(*jsPlugin) {
	return func(js *jsPlugin) {
		js.classExtractor = pce
	}
}

// WithMetadataExtractor option replaces the default JSMetadataExtractor function of a JSPlugin.
func WithMetadataExtractor(pme JSMetadataExtractor) func(*jsPlugin) {
	return func(js *jsPlugin) {
		js.metadataExtractor = pme
	}
}

// JSOption is an option that overrides a default value of a JSPlugin.
type JSOption func(*jsPlugin)

// JSPlugin interface represents a JavaScript plugin.
type JSPlugin interface {
	Close()
	PluginPath() string
	Metadata() *Metadata

	Navigation(ctx context.Context) (navigation.Navigation, error)
	Register(ctx context.Context, dashboardAPIAddress string) (Metadata, error)
	Print(ctx context.Context, object runtime.Object) (PrintResponse, error)
	PrintTab(ctx context.Context, object runtime.Object) (TabResponse, error)
	ObjectStatus(ctx context.Context, object runtime.Object) (ObjectStatusResponse, error)
	HandleAction(ctx context.Context, actionName string, payload action.Payload) error
	Content(ctx context.Context, contentPath string) (component.ContentResponse, error)
}

type jsPlugin struct {
	loop *eventloop.EventLoop

	metadata    *Metadata
	pluginClass *goja.Object
	pluginPath  string

	runtimeFactory    JSRuntimeFactory
	classExtractor    JSClassExtractor
	metadataExtractor JSMetadataExtractor

	mu     sync.Mutex
	ctx    context.Context
	logger log.Logger
}

var _ JSPlugin = (*jsPlugin)(nil)

// NewJSPlugin creates a new instances of a JavaScript plugin.
func NewJSPlugin(ctx context.Context, pluginPath string, dashboardClientFactory octant.DashboardClientFactory, options ...JSOption) (*jsPlugin, error) {
	plugin := &jsPlugin{
		ctx:               ctx,
		pluginPath:        pluginPath,
		runtimeFactory:    javascript.CreateRuntimeLoop,
		classExtractor:    javascript.ExtractDefaultClass,
		metadataExtractor: extractMetadata,
	}

	for _, o := range options {
		o(plugin)
	}

	loop, err := plugin.runtimeFactory(ctx, pluginPath)

	if err != nil {
		return nil, fmt.Errorf("initializing runtime: %w", err)
	}

	var pluginClass *goja.Object
	var metadata *Metadata

	errCh := make(chan error)

	loop.RunOnLoop(func(vm *goja.Runtime) {
		var err error

		buf, err := ioutil.ReadFile(pluginPath)
		if err != nil {
			errCh <- fmt.Errorf("reading script: %w", err)
		}
		program, err := goja.Compile(pluginPath, string(buf), false)
		if err != nil {
			errCh <- fmt.Errorf("compiling: %w", err)
		}
		_, err = vm.RunProgram(program)
		if err != nil {
			errCh <- fmt.Errorf("script execution: %w", err)
		}

		// Convert these to use require.RegisterNativeModule
		vm.Set("httpClient", javascript.CreateHTTPClientObject(vm, pluginClass))
		vm.Set("dashboardClient", dashboardClientFactory.Create(ctx, vm))

		pluginClass, err = plugin.classExtractor(vm)
		if err != nil {
			errCh <- fmt.Errorf("loading pluginClass: %w", err)
		}

		metadata, err = plugin.metadataExtractor(vm, pluginClass)
		if err != nil {
			errCh <- fmt.Errorf("loading metadata: %w", err)
		}

		errCh <- nil

	})

	err = <-errCh
	if err != nil {
		return nil, fmt.Errorf("javascript loop: %w", err)
	}

	plugin.loop = loop
	plugin.pluginClass = pluginClass
	plugin.metadata = metadata

	return plugin, nil
}

// Close closes the dashboard client connection.
func (t *jsPlugin) Close() {
	t.loop.Stop()
}

// PluginPath returns the pluginPath.
func (t *jsPlugin) PluginPath() string {
	return t.pluginPath
}

// Navigation returns the navigation for a JavaScript plugin.
func (t *jsPlugin) Navigation(_ context.Context) (navigation.Navigation, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	nav := navigation.Navigation{}
	errCh := make(chan error)

	t.loop.RunOnLoop(func(vm *goja.Runtime) {
		handler, err := vm.RunString("_concretePlugin.navigationHandler")
		if err != nil {
			errCh <- fmt.Errorf("unable to load navigationHandler from plugin: %w", err)
			return
		}

		cHandler, ok := goja.AssertFunction(handler)
		if !ok {
			errCh <- fmt.Errorf("navigationHandler is not callable")
			return
		}

		s, err := cHandler(t.pluginClass)
		if err != nil {
			errCh <- fmt.Errorf("calling navigationHandler: %w", err)
			return
		}

		jsonNav, err := json.Marshal(s.Export())
		if err != nil {
			errCh <- fmt.Errorf("unable to marshal navigation json: %w", err)
			return
		}

		if err := json.Unmarshal(jsonNav, &nav); err != nil {
			errCh <- fmt.Errorf("unable to unmarshal navigation json: %w", err)
			return
		}
		errCh <- nil
	})

	err := <-errCh
	if err != nil {
		return nav, err
	}

	return nav, nil
}

// Content returns the content response for a JavaScript plugin acting as a module.
func (t *jsPlugin) Content(ctx context.Context, contentPath string) (component.ContentResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	cr := component.ContentResponse{}
	errCh := make(chan error)

	t.loop.RunOnLoop(func(vm *goja.Runtime) {
		clientID := ocontext.WebsocketClientIDFrom(ctx)

		handler, err := vm.RunString("_concretePlugin.contentHandler")
		if err != nil {
			errCh <- fmt.Errorf("unable to load contentHandler from plugin: %w", err)
			return
		}

		cHandler, ok := goja.AssertFunction(handler)
		if !ok {
			errCh <- fmt.Errorf("contentHandler is not callable")
			return
		}
		obj := vm.NewObject()
		if err := obj.Set("contentPath", vm.ToValue(contentPath)); err != nil {
			errCh <- fmt.Errorf("unable to set contentPath: %w", err)
			return
		}
		if err := obj.Set("clientID", vm.ToValue(clientID)); err != nil {
			errCh <- fmt.Errorf("unable to set clientID: %w", err)
		}
		s, err := cHandler(t.pluginClass, obj)
		if err != nil {
			errCh <- fmt.Errorf("calling contentHandler: %w", err)
			return
		}

		pluginResp := s.ToObject(vm)
		if pluginResp == nil {
			errCh <- fmt.Errorf("empty contentResponse")
			return
		}

		content := pluginResp.Get("content")
		if content == goja.Undefined() {
			errCh <- fmt.Errorf("unable to get content from contentResponse")
			return
		}

		contentObj, ok := content.Export().(map[string]interface{})
		if !ok {
			errCh <- fmt.Errorf("unable to get content as map from contentResponse")
			return
		}

		rawTitle, ok := contentObj["title"]
		if ok {
			titles, ok := rawTitle.([]interface{})
			if !ok {
				errCh <- fmt.Errorf("unable to get title array from content")
				return
			}
			for i, c := range titles {
				realTitle, err := javascript.ConvertToComponent(fmt.Sprintf("title[%d]", i), c)
				if err != nil {
					errCh <- fmt.Errorf("unable to extract title: %w", err)
					return
				}

				title, ok := realTitle.(component.TitleComponent)
				if !ok {
					errCh <- fmt.Errorf("unable to convert component to TitleComponent")
					return
				}
				cr.Title = append(cr.Title, title)
			}
		}

		rawComponents, ok := contentObj["viewComponents"]
		if !ok {
			errCh <- fmt.Errorf("unable to get viewComponents from content")
			return
		}

		components, ok := rawComponents.([]interface{})
		if !ok {
			errCh <- fmt.Errorf("unable to get viewComponents list")
			return
		}

		for i, c := range components {
			realComponent, err := javascript.ConvertToComponent(fmt.Sprintf("viewComponent[%d]", i), c)
			if err != nil {
				errCh <- fmt.Errorf("unable to extract component: %w", err)
				return
			}
			cr.Add(realComponent)
		}

		rawButtonGroup, ok := contentObj["buttonGroup"]
		if ok {
			realButtonGroup, err := javascript.ConvertToComponent("buttonGroup", rawButtonGroup)
			if err != nil {
				errCh <- fmt.Errorf("unable to extract buttonGroup: %w", err)
				return
			}

			buttonGroup, ok := realButtonGroup.(*component.ButtonGroup)
			if !ok {
				errCh <- fmt.Errorf("unable to convert extracted component to buttonGroup")
				return
			}

			cr.ButtonGroup = buttonGroup
		}
		errCh <- nil
	})

	if err := <-errCh; err != nil {
		return cr, err
	}
	return cr, nil
}

// Metadata returns the JavaScript plugins metadata.
func (t *jsPlugin) Metadata() *Metadata {
	return t.metadata
}

// Register is not implemented for JavaScript plugins.
func (t *jsPlugin) Register(_ context.Context, _ string) (Metadata, error) {
	return Metadata{}, fmt.Errorf("not implemented")
}

// PrintTab returns the tab response from a JavaScript plugins tab handler.
func (t *jsPlugin) PrintTab(ctx context.Context, object runtime.Object) (TabResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	tabResponse, err := t.objectRequestCall(ctx, "tabHandler", object)
	if err != nil {
		return TabResponse{}, err
	}

	tab := tabResponse.Get("tab")
	if tab == goja.Undefined() {
		return TabResponse{}, fmt.Errorf("tab property not found")
	}

	contents, ok := tab.Export().(map[string]interface{})
	if !ok {
		return TabResponse{}, fmt.Errorf("unable to export tab contents")
	}

	cTab := &component.Tab{}
	if name, ok := contents["name"]; ok {
		cTab.Contents = *component.NewFlexLayout(name.(string))
		cTab.Name = name.(string)
	}

	if contents, ok := contents["contents"]; ok {
		realComponent, err := javascript.ConvertToComponent("tab contents", contents)
		if err != nil {
			return TabResponse{}, fmt.Errorf("unable to extract component: %w", err)
		}
		realFlexLayout := *realComponent.(*component.FlexLayout)
		cTab.Contents.Config = realFlexLayout.Config
	}

	return TabResponse{
		Tab: cTab,
	}, nil
}

// ObjectStats returns the object status from a JavaScript plugins object status handler.
func (t *jsPlugin) ObjectStatus(ctx context.Context, object runtime.Object) (ObjectStatusResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	osResponse, err := t.objectRequestCall(ctx, "objectStatusHandler", object)
	if err != nil {
		return ObjectStatusResponse{}, err
	}

	objStatus := osResponse.Get("objectStatus")
	if objStatus == goja.Undefined() {
		return ObjectStatusResponse{}, fmt.Errorf("objectStatus property not found")
	}

	mapObjStatus, ok := objStatus.Export().(map[string]interface{})
	if !ok {
		return ObjectStatusResponse{}, fmt.Errorf("unable to get objectStatus map")
	}

	jsonOS, err := json.Marshal(mapObjStatus)
	if err != nil {
		return ObjectStatusResponse{}, fmt.Errorf("unable to marshal podSummary: %w", err)
	}

	var podSummary component.PodSummary
	if err := json.Unmarshal(jsonOS, &podSummary); err != nil {
		return ObjectStatusResponse{}, fmt.Errorf("unable to unmarshal podSummary: %w", err)
	}

	return ObjectStatusResponse{
		ObjectStatus: podSummary,
	}, nil
}

// HandleAction calls the JavaScript plugins action handler.
func (t *jsPlugin) HandleAction(ctx context.Context, actionPath string, payload action.Payload) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	errCh := make(chan error)

	t.loop.RunOnLoop(func(vm *goja.Runtime) {
		clientID := ocontext.WebsocketClientIDFrom(ctx)

		handler, err := vm.RunString("_concretePlugin.actionHandler")
		if err != nil {
			errCh <- fmt.Errorf("unable to load actionHandler from plugin: %w", err)
			return
		}

		cHandler, ok := goja.AssertFunction(handler)
		if !ok {
			errCh <- fmt.Errorf("actionHandler is not callable")
			return
		}

		var pl map[string]interface{}
		pl = payload

		obj := vm.NewObject()
		if err := obj.Set("actionName", vm.ToValue(actionPath)); err != nil {
			errCh <- fmt.Errorf("unable to set actionName: %w", err)
			return
		}
		if err := obj.Set("payload", pl); err != nil {
			errCh <- fmt.Errorf("unable to set payload: %w", err)
			return
		}
		if err := obj.Set("clientID", clientID); err != nil {
			errCh <- fmt.Errorf("unable to set clientID: %w", err)
			return
		}

		s, err := cHandler(t.pluginClass, obj)
		if err != nil {
			errCh <- fmt.Errorf("calling actionHandler: %w", err)
			return
		}

		if s != goja.Undefined() {
			if jsErr := s.ToObject(vm); jsErr != nil {
				errStr := jsErr.Get("error")
				if errStr != goja.Undefined() {
					errCh <- fmt.Errorf("%s actionHandler: %q", t.pluginPath, jsErr.Get("error"))
					return
				}
			}
		}
		errCh <- nil
	})

	if err := <-errCh; err != nil {
		return err
	}

	return nil
}

// Print returns the print response from the JavaScript plugins print handler.
func (t *jsPlugin) Print(ctx context.Context, object runtime.Object) (PrintResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	printResponse, err := t.objectRequestCall(ctx, "printHandler", object)
	if err != nil {
		return PrintResponse{}, err
	}

	sections, ok := printResponse.Export().(map[string]interface{})
	if !ok {
		return PrintResponse{}, fmt.Errorf("unable to parse printHandler response sections")
	}

	var configSections []component.SummarySection
	var statusSections []component.SummarySection
	var flexItems []component.FlexLayoutItem

	for k, v := range sections {
		switch k {
		case "config":
			ss, err := javascript.ConvertToSections(k, v)
			if err != nil {
				return PrintResponse{}, fmt.Errorf("error extracting sections: %w", err)
			}
			configSections = append(configSections, ss...)
		case "status":
			ss, err := javascript.ConvertToSections(k, v)
			if err != nil {
				return PrintResponse{}, fmt.Errorf("error extracting sections: %w", err)
			}
			statusSections = append(statusSections, ss...)
		case "items":
			ss, err := javascript.ConvertToItems(k, v)
			if err != nil {
				return PrintResponse{}, fmt.Errorf("error extracting items: %w", err)
			}
			flexItems = append(flexItems, ss...)
		default:
			return PrintResponse{}, fmt.Errorf("unknown printHandler response section: %s", k)
		}
	}

	var response PrintResponse
	response.Config = configSections
	response.Status = statusSections
	response.Items = flexItems

	return response, nil
}

func (t *jsPlugin) objectRequestCall(ctx context.Context, handlerName string, object runtime.Object) (*goja.Object, error) {
	errCh := make(chan error)
	var response *goja.Object

	t.loop.RunOnLoop(func(vm *goja.Runtime) {
		clientID := ocontext.WebsocketClientIDFrom(ctx)

		handler, err := vm.RunString(fmt.Sprintf("_concretePlugin.%s", handlerName))
		if err != nil {
			errCh <- fmt.Errorf("unable to load %s from plugin: %w", handlerName, err)
			return
		}

		cHandler, ok := goja.AssertFunction(handler)
		if !ok {
			errCh <- fmt.Errorf("%s is not callable", handlerName)
			return
		}

		obj := vm.NewObject()
		if err := obj.Set("object", vm.ToValue(object)); err != nil {
			errCh <- fmt.Errorf("unable to set object: %w", err)
			return
		}
		if err := obj.Set("clientID", vm.ToValue(clientID)); err != nil {
			errCh <- fmt.Errorf("unable to set clientID: %w", err)
			return
		}
		s, err := cHandler(t.pluginClass, obj)
		if err != nil {
			errCh <- err
			return
		}

		response = s.ToObject(vm)
		if response == nil {
			errCh <- fmt.Errorf("no status found")
			return
		}
		errCh <- nil
	})

	if err := <-errCh; err != nil {
		return nil, err
	}

	return response, nil
}

func extractMetadata(vm *goja.Runtime, pluginValue goja.Value) (*Metadata, error) {
	metadata := new(Metadata)

	this := pluginValue.ToObject(vm)
	if this == nil {
		return nil, fmt.Errorf("unable to construct `this` from plugin class")
	}

	metadata.Name = this.Get("name").String()
	if metadata.Name == "" {
		return nil, fmt.Errorf("name is a required property")
	}

	metadata.Description = this.Get("description").String()
	if metadata.Description == "" {
		return nil, fmt.Errorf("description is a required property")
	}

	metadata.Capabilities.IsModule = this.Get("isModule").ToBoolean()

	if capability, ok := this.Get("capabilities").Export().(map[string]interface{}); ok {
		for k, v := range capability {
			switch k {
			case "supportPrinterConfig":
				GVKs, err := javascript.ConvertToGVKs(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsPrinterConfig = append(metadata.Capabilities.SupportsPrinterConfig, GVKs...)
			case "supportPrinterStatus":
				GVKs, err := javascript.ConvertToGVKs(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsPrinterStatus = append(metadata.Capabilities.SupportsPrinterStatus, GVKs...)
			case "supportPrinterItems":
				GVKs, err := javascript.ConvertToGVKs(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsPrinterItems = append(metadata.Capabilities.SupportsPrinterItems, GVKs...)
			case "supportObjectStatus":
				GVKs, err := javascript.ConvertToGVKs(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsObjectStatus = append(metadata.Capabilities.SupportsObjectStatus, GVKs...)
			case "supportTab":
				GVKs, err := javascript.ConvertToGVKs(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsTab = append(metadata.Capabilities.SupportsTab, GVKs...)
			case "actionNames":
				actions, err := javascript.ConvertToActions(v)
				if err != nil {
					return nil, fmt.Errorf("extractActions: %w", err)
				}
				metadata.Capabilities.ActionNames = append(metadata.Capabilities.ActionNames, actions...)
			default:
				fmt.Printf("unknown capabilitiy: %s\n", k)
			}
		}
	} else {
		return nil, fmt.Errorf("unable to get capabilites for plugin class")
	}

	return metadata, nil
}
