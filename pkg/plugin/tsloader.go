package plugin

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dop251/goja"
)

type TSLoader struct {
	vm *goja.Runtime
}

func NewTSLoader(vm *goja.Runtime) (*TSLoader, error) {
	box, err := rice.FindBox("_files")
	if err != nil {
		return nil, err
	}

	if err := initTypescriptServices(box, vm); err != nil {
		return nil, err
	}
	return &TSLoader{
		vm: vm,
	}, nil
}

func (t *TSLoader) SourceLoader(path string) ([]byte, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return t.loadFile(path)
	}

	if fi.IsDir() {
		return t.loadDir(path)
	}
	return t.loadFile(path)
}

func (t *TSLoader) loadDir(path string) ([]byte, error) {
	var ret []byte

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		fp := filepath.Join(path, fi.Name())
		if fi.IsDir() {
			v, err := t.loadDir(fp)
			if err != nil {
				return nil, err
			}
			ret = append(ret[:], v[:]...)
		} else {
			if IsJavaScriptPlugin(fp) {
				v, err := t.loadFile(fp)
				if err != nil {
					return nil, err
				}
				ret = append(ret[:], v[:]...)
			}
			continue
		}
	}
	return ret, nil
}

func (t *TSLoader) loadFile(path string) ([]byte, error) {
	fmt.Println(path)
	fileName := filepath.Base(path)
	swd := filepath.Dir(path)

	fis, err := ioutil.ReadDir(swd)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		if !fi.IsDir() && strings.Contains(fi.Name(), fileName) {
			if IsTypescriptPlugin(fi.Name()) {
				fmt.Printf("transpiling: %s\n", fi.Name())
				return t.TranspileModule(filepath.Join(swd, fi.Name()))
			}
			return ioutil.ReadFile(filepath.Join(swd, fi.Name()))
		}
	}
	return nil, fmt.Errorf("No matching file found for %s", path)
}

func (t *TSLoader) TranspileModule(path string) ([]byte, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	scriptValue := t.vm.ToValue(string(buf))
	t.vm.Set("scriptValue", scriptValue)

	v, err := t.vm.RunString(`ts.transpile(scriptValue, compilerOptions, undefined, undefined, "plugin");`)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf("no value from transpileModule for %s", path)
	}

	return []byte(v.String()), err
}

func (t *TSLoader) TranspileAndRun(path string) (goja.Value, error) {
	b, err := t.TranspileModule(path)

	if err != nil {
		return goja.Undefined(), err
	}
	return t.vm.RunString(string(b[:]))
}

func initTypescriptServices(box *rice.Box, vm *goja.Runtime) error {
	f, err := box.Open("typescriptServices.js")
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = vm.RunString(string(buf))
	if err != nil {
		return err
	}

	ts := vm.Get("ts")
	if ts == nil {
		return fmt.Errorf("unable to load typescriptServices module")
	}
	tp := ts.ToObject(vm).Get("transpile")

	_, ok := goja.AssertFunction(tp)
	if !ok {
		return fmt.Errorf("no transpile function found")
	}

	vm.Set("global", vm.GlobalObject())
	vm.Set("self", vm.GlobalObject())

	_, err = vm.RunString(`
var module = { exports: {} };
var exports = module.exports;
var compilerOptions = { module: ts.ModuleKind.SystemJS };
`)

	if err != nil {
		return err
	}
	return nil
}
