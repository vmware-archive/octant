package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
	"sync"

	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/navigation"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/vmware-tanzu/octant/pkg/plugin/api"
	"github.com/vmware-tanzu/octant/pkg/plugin/console"
	"github.com/vmware-tanzu/octant/pkg/store"
	"k8s.io/apimachinery/pkg/runtime"
)

func IsTypescriptPlugin(pluginName string) bool {
	return strings.Contains(pluginName, ".ts")
}

func IsJavaScriptPlugin(pluginName string) bool {
	return IsTypescriptPlugin(pluginName) || strings.Contains(pluginName, ".js")
}

type dashboardClient struct {
	client *api.Client
	vm     *goja.Runtime
}

func (d *dashboardClient) Get(c goja.FunctionCall) goja.Value {
	var key store.Key
	obj := c.Argument(0).ToObject(d.vm)
	d.vm.ExportTo(obj, &key)

	u, err := d.client.Get(context.TODO(), key)
	if err != nil {
		return d.vm.NewGoError(err)
	}

	return d.vm.ToValue(u.Object)
}

func (d *dashboardClient) List(c goja.FunctionCall) goja.Value {
	var key store.Key
	obj := c.Argument(0).ToObject(d.vm)
	d.vm.ExportTo(obj, &key)

	u, err := d.client.List(context.TODO(), key)
	if err != nil {
		return d.vm.NewGoError(err)
	}

	return d.vm.ToValue(u.Object)
}

func (d *dashboardClient) Create(c goja.FunctionCall) goja.Value {
	var u *unstructured.Unstructured
	obj := c.Argument(0).ToObject(d.vm)
	d.vm.ExportTo(obj, &u)

	err := d.client.Create(context.TODO(), u)
	if err != nil {
		return d.vm.NewGoError(err)
	}

	return goja.Undefined()
}

func (d *dashboardClient) Update(c goja.FunctionCall) goja.Value {
	var u *unstructured.Unstructured
	obj := c.Argument(0).ToObject(d.vm)
	d.vm.ExportTo(obj, &u)

	err := d.client.Update(context.TODO(), u)
	if err != nil {
		return d.vm.NewGoError(err)
	}

	return goja.Undefined()
}

func createClientObject(d *dashboardClient) goja.Value {
	obj := d.vm.NewObject()
	obj.Set("Get", d.Get)
	obj.Set("List", d.List)
	obj.Set("Create", d.Create)
	obj.Set("Update", d.Update)
	return obj
}

type pluginRuntimeFactory func(context.Context, string, bool) (*goja.Runtime, *TSLoader, error)
type pluginClassExtractor func(*goja.Runtime) (*goja.Object, error)
type pluginMetadataExtractor func(*goja.Runtime, goja.Value) (*Metadata, error)

func NewJSPlugin(ctx context.Context, client *api.Client, pluginPath string, prf pluginRuntimeFactory, pce pluginClassExtractor, pme pluginMetadataExtractor) (*jsPlugin, error) {
	vm, tsl, err := prf(ctx, pluginPath, IsTypescriptPlugin(pluginPath))

	if err != nil {
		return nil, fmt.Errorf("initializing runtime: %w", err)
	}

	if tsl != nil {
		_, err = tsl.TranspileAndRun(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("script transpile and execution: %w", err)
		}
	} else {
		buf, err := ioutil.ReadFile(pluginPath)
		if err != nil {
			return nil, fmt.Errorf("reading script: %w", err)
		}
		_, err = vm.RunString(string(buf[:]))
		if err != nil {
			return nil, fmt.Errorf("script execution: %w", err)
		}
	}

	pluginClass, err := pce(vm)
	if err != nil {
		return nil, fmt.Errorf("loading pluginClass: %w", err)
	}

	metadata, err := pme(vm, pluginClass)
	if err != nil {
		return nil, fmt.Errorf("loading metadata: %w", err)
	}

	plugin := &jsPlugin{
		vm:                vm,
		metadata:          metadata,
		client:            client,
		runtimeFactory:    prf,
		classExtractor:    pce,
		metadataExtractor: pme,
		pluginPath:        pluginPath,
		ctx:               ctx,
	}

	return plugin, nil
}

type jsPlugin struct {
	vm *goja.Runtime

	metadata    *Metadata
	pluginPath  string

	runtimeFactory    pluginRuntimeFactory
	classExtractor    pluginClassExtractor
	metadataExtractor pluginMetadataExtractor

	client  *api.Client
	mu      sync.Mutex
	ctx     context.Context
}

func (t *jsPlugin) Reload() error {
	fmt.Println("calling reload")

	logger := log.From(t.ctx)
	logger.Infof("reloading plugin: %s\n", t.pluginPath)
	t.mu.Lock()
	defer t.mu.Unlock()

	vm, tsl, err := t.runtimeFactory(t.ctx, t.pluginPath, IsTypescriptPlugin(t.pluginPath))
	if err != nil {
		return fmt.Errorf("realod: runtime initialization: %w", err)
	}

	if tsl != nil {
		_, err = tsl.TranspileAndRun(t.pluginPath)
		if err != nil {
			return fmt.Errorf("reload: transpile and execute: %w", err)
		}
	} else {
		buf, err := ioutil.ReadFile(t.pluginPath)
		if err != nil {
			return fmt.Errorf("reload: reading script: %w", err)
		}
		_, err = vm.RunString(string(buf[:]))
		if err != nil {
			return fmt.Errorf("reload: script execution: %w", err)
		}
	}

	pluginClass, err := t.classExtractor(vm)
	if err != nil {
		return err
	}

	metadata, err := t.metadataExtractor(vm, pluginClass)
	if err != nil {
		return err
	}

	t.vm = vm
	t.metadata = metadata

	return nil
}

func (t *jsPlugin) Navigation(ctx context.Context) (navigation.Navigation, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	nav := navigation.Navigation{}

	handler, err := t.vm.RunString("_concretePlugin.navigationHandler")
	if err != nil {
		return nav, fmt.Errorf("unable to load navigationHandler from plugin: %w", err)
	}

	cHandler, ok := goja.AssertFunction(handler)
	if !ok {
		return nav, fmt.Errorf("navigationHandler is not callable")
	}

	s, err := cHandler(goja.Undefined())
	if err != nil {
		return nav, fmt.Errorf("calling navigationHandler: %w", err)
	}

	jsonNav, err := json.Marshal(s.Export())
	if err != nil {
		return nav, fmt.Errorf("unable to marshal navigation json: %w", err)
	}

	if err := json.Unmarshal(jsonNav, &nav); err != nil {
		return nav, fmt.Errorf("unable to unmarshal navigation json: %w", err)
	}

	return nav, nil
}

func (t *jsPlugin) Content(ctx context.Context, contentPath string) (component.ContentResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	cr := component.ContentResponse{}

	handler, err := t.vm.RunString("_concretePlugin.contentHandler")
	if err != nil {
		return cr, fmt.Errorf("unable to load contentHandler from plugin: %w", err)
	}

	cHandler, ok := goja.AssertFunction(handler)
	if !ok {
		return cr, fmt.Errorf("contentHandler is not callable")
	}

	s, err := cHandler(goja.Undefined(), t.vm.ToValue(contentPath))
	if err != nil {
		return cr, fmt.Errorf("calling contentHandler: %w", err)
	}

	jsonCr, err := json.Marshal(s.Export())
	if err != nil {
		return cr, fmt.Errorf("unable to marshal content json: %w", err)
	}

	if err := json.Unmarshal(jsonCr, &cr); err != nil {
		return cr, fmt.Errorf("unable to unmarshal content json: %w", err)
	}

	return cr, nil
}

func (t *jsPlugin) Metadata() *Metadata {
	return t.metadata
}

func (t *jsPlugin) Register(ctx context.Context, dashboardAPIAddress string) (Metadata, error) {
	return Metadata{}, fmt.Errorf("not implemented")
}

func (t *jsPlugin) PrintTab(ctx context.Context, object runtime.Object) (TabResponse, error) {
	return TabResponse{}, fmt.Errorf("not implemented")
}

func (t *jsPlugin) ObjectStatus(ctx context.Context, object runtime.Object) (ObjectStatusResponse, error) {
	return ObjectStatusResponse{}, fmt.Errorf("not implemented")
}

func (t *jsPlugin) HandleAction(ctx context.Context, actionName string, payload action.Payload) error {
	return fmt.Errorf("not implemented")
}

func (t *jsPlugin) Print(_ context.Context, object runtime.Object) (PrintResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	handler, err := t.vm.RunString("_concretePlugin.printHandler")
	if err != nil {
		return PrintResponse{}, fmt.Errorf("unable to load printHandler from plugin: %w", err)
	}

	cHandler, ok := goja.AssertFunction(handler)
	if !ok {
		return PrintResponse{}, fmt.Errorf("printHandler is not callable")
	}

	gc := dashboardClient{
		client: t.client,
		vm:     t.vm,
	}

	key, err := store.KeyFromObject(object)
	if err != nil {
		return PrintResponse{}, fmt.Errorf("ts plugin generate key: %w", err)
	}

	obj := t.vm.NewObject()
	obj.Set("client", createClientObject(&gc))
	obj.Set("objectKey", t.vm.ToValue(&key))

	s, err := cHandler(goja.Undefined(), obj)
	if err != nil {
		return PrintResponse{}, err
	}

	printResponse := s.ToObject(t.vm)
	if printResponse == nil {
		return PrintResponse{}, fmt.Errorf("no status found")
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
			ss, err := extractSections(k, v)
			if err != nil {
				return PrintResponse{}, fmt.Errorf("error extracting sections: %w", err)
			}
			configSections = append(configSections, ss...)
		case "status":
			ss, err := extractSections(k, v)
			if err != nil {
				return PrintResponse{}, fmt.Errorf("error extracting sections: %w", err)
			}
			statusSections = append(statusSections, ss...)
		case "items":
			ss, err := extractItems(k, v)
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

func extractItems(name string, i interface{}) ([]component.FlexLayoutItem, error) {
	var items []component.FlexLayoutItem

	v, ok := i.([]interface{})
	if !ok {
		return items, fmt.Errorf("unable to parse printHandler %s summary items", name)
	}

	for idx, ii := range v {
		mapItem, ok := ii.(map[string]interface{})
		if !ok {
			return items, fmt.Errorf("unable to parse %s summary items", name)
		}
		flexLayoutItem := component.FlexLayoutItem{}
		jsonSS, err := json.Marshal(mapItem)
		if err != nil {
			return items, fmt.Errorf("unable to marshal json in position %d in %s", idx, name)
		}
		if err := json.Unmarshal(jsonSS, &flexLayoutItem); err != nil {
			return items, fmt.Errorf("unable to unmarshal json in position %d in %s", idx, name)
		}
		items = append(items, flexLayoutItem)
	}

	return items, nil
}

func extractSections(name string, i interface{}) ([]component.SummarySection, error) {
	var sections []component.SummarySection

	v, ok := i.([]interface{})
	if !ok {
		return sections, fmt.Errorf("unable to parse printHandler %s summary sections", name)
	}

	for idx, ii := range v {
		mapSummarySection, ok := ii.(map[string]interface{})
		if !ok {
			return sections, fmt.Errorf("unable to parse %s summary section", name)
		}
		ss := component.SummarySection{}
		jsonSS, err := json.Marshal(mapSummarySection)
		if err != nil {
			return sections, fmt.Errorf("unable to marshal json GVK in position %d in %s", idx, name)
		}
		if err := json.Unmarshal(jsonSS, &ss); err != nil {
			return sections, fmt.Errorf("unable to unmarshal json GVK in position %d in %s", idx, name)
		}
		sections = append(sections, ss)
	}

	return sections, nil
}

func ExtractDefaultClass(vm *goja.Runtime) (*goja.Object, error) {
	// This is the location of a export default class that implements the Octant
	// TypeScript module definition.
	instantiateClass := "var _concretePlugin = new module.exports.default; _concretePlugin"
	// This is the library name the Octant webpack configuration uses.
	if vm.Get("_octantPlugin") != nil {
		instantiateClass = "var _concretePlugin = new _octantPlugin; _concretePlugin"
	}

	v, err := vm.RunString(instantiateClass)
	if err != nil {
		return nil, fmt.Errorf("unable to create new plugin: %w", err)
	}
	pluginClass := v.ToObject(vm)
	return pluginClass, nil
}

func extractGvk(name string, i interface{}) ([]schema.GroupVersionKind, error) {
	gvks, ok := i.([]interface{})
	if !ok {
		fmt.Errorf("unable to parse GVK list for supportPrinterConfig")
	}
	var gvkList []schema.GroupVersionKind
	for i, ii := range gvks {
		mapGvk, ok := ii.(map[string]interface{})
		if !ok {
			return gvkList, fmt.Errorf("unable to parse GVK in position %d in supportPrinterConfig", i)
		}
		gvk := schema.GroupVersionKind{}

		jsonGvk, err := json.Marshal(mapGvk)
		if err != nil {
			return gvkList, fmt.Errorf("unable to marshal json GVK in position %d in supportPrinterConfig", i)
		}

		if err := json.Unmarshal(jsonGvk, &gvk); err != nil {
			return gvkList, fmt.Errorf("unablet to unmarshal json GVK in position %d in supportPrinterConfig", i)
		}

		gvkList = append(gvkList, gvk)
	}
	return gvkList, nil
}

func ExtractMetadata(vm *goja.Runtime, pluginValue goja.Value) (*Metadata, error) {
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

	if capab, ok := this.Get("capabilities").Export().(map[string]interface{}); ok {
		for k, v := range capab {
			switch k {
			case "supportPrinterConfig":
				gvks, err := extractGvk(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsPrinterConfig = append(metadata.Capabilities.SupportsPrinterConfig, gvks...)
			case "supportPrinterStatus":
				gvks, err := extractGvk(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsPrinterStatus = append(metadata.Capabilities.SupportsPrinterStatus, gvks...)
			case "supportPrinterItems":
				gvks, err := extractGvk(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsPrinterItems = append(metadata.Capabilities.SupportsPrinterItems, gvks...)
			case "supportObjectStatus":
				gvks, err := extractGvk(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsObjectStatus = append(metadata.Capabilities.SupportsObjectStatus, gvks...)
			case "supportTab":
				gvks, err := extractGvk(k, v)
				if err != nil {
					return nil, fmt.Errorf("extractGvks: %w", err)
				}
				metadata.Capabilities.SupportsTab = append(metadata.Capabilities.SupportsTab, gvks...)
			default:
				fmt.Printf("unknown capabilitiy: %s\n", k)
			}
		}
	} else {
		return nil, fmt.Errorf("unable to get capabilites for plugin class")
	}

	return metadata, nil
}

func CreateRuntime(ctx context.Context, logName string, typescript bool) (*goja.Runtime, *TSLoader, error) {
	vm := goja.New()

	vm.Set("global", vm.GlobalObject())
	vm.Set("self", vm.GlobalObject())

	_, err := vm.RunString(`
var module = { exports: {} };
var exports = module.exports;
`)
	if err != nil {
		return nil, nil, fmt.Errorf("runtime global values: %w", err)
	}

	var tsl *TSLoader
	if typescript {
		tsl, err = NewTSLoader(vm)
		if err != nil {
			return nil, nil, err
		}
	}

	registry := require.NewRegistryWithLoader(tsl.SourceLoader)
	registry.Enable(vm)

	logger := log.From(ctx).Named(logName)
	console.CustomInit(logger)
	console.Enable(vm)

	// This maps Caps fields to lower fields from struct to Object based on the JSON annotations.
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
	return vm, tsl, nil
}
