/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package tsgen

import (
	"bytes"
	"embed"
	"encoding/gob"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/iancoleman/strcase"
)

//go:embed ts-support
var tsSupportFiles embed.FS

var tsSupportDir string = "ts-support"

// TSGen generates typescript for components.
type TSGen struct {
	disableFormat bool
}

// NewTSGen creates an instance of TSGen.
func NewTSGen(options ...Option) (*TSGen, error) {
	opts := makeDefaultOptions(options...)

	t := &TSGen{
		disableFormat: opts.disableFormat,
	}
	return t, nil
}

// Names returns all the component names ina
func (tg *TSGen) Names(source string) ([]string, error) {
	fileSet := token.NewFileSet()

	// TODO: this may not be needed as we can glean this from source.
	dirs := []string{
		filepath.Join("pkg", "view", "component"),
	}

	var files []*ast.File
	for _, dir := range dirs {
		wd := filepath.Join(".", dir)

		pkg, err := build.ImportDir(wd, build.ImportComment)
		if err != nil {
			return nil, fmt.Errorf("import dir %s: %w", wd, err)
		}

		for _, goFile := range pkg.GoFiles {
			fileName := filepath.Join(pkg.Dir, goFile)
			file, err := parser.ParseFile(fileSet, fileName, nil, parser.ParseComments|parser.AllErrors)
			if err != nil {
				return nil, fmt.Errorf("parsing %s: %w", fileName, err)
			}
			files = append(files, file)
		}
	}

	names, err := componentsInPackage(fileSet, files, source)
	if err != nil {
		return nil, fmt.Errorf("find components in package: %w", err)
	}

	return names, nil
}

// Reflect reflects on list of components. It generates a temporary binary and runs that to
// get type information for components.
func (tg *TSGen) Reflect(names []string) (*Model, error) {
	dir, err := ioutil.TempDir(".", "reflect")
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := os.RemoveAll(dir); cErr != nil {
			log.Fatalf("unable to remove temporary directory")
		}
	}()

	b, err := tg.ReflectTemplate(names)
	if err != nil {
		return nil, fmt.Errorf("create reflect program")
	}

	reflectProgram := filepath.Join(dir, "reflect.go")

	if err := ioutil.WriteFile(reflectProgram, b, 0666); err != nil {
		return nil, fmt.Errorf("create reflect program: %w", err)
	}

	cmd := exec.Command("go", "run", reflectProgram)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf(stderr.String())
		return nil, fmt.Errorf("run reflect program: %w", err)
	}

	m := &Model{}
	if err := gob.NewDecoder(&stdout).Decode(m); err != nil {
		return nil, fmt.Errorf("decode reflect model: %w", err)
	}

	return m, nil
}

// ReflectTemplate generates a reflect template.
func (tg *TSGen) ReflectTemplate(names []string) ([]byte, error) {
	b, err := tsSupportFiles.ReadFile(path.Join(tsSupportDir, "reflect.go.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("load reflect template: %w", err)
	}

	t, err := template.New("reflect").Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("parse reflect template: %w", err)
	}

	opts := reflectOptions{
		Names: names,
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, opts); err != nil {
		return nil, fmt.Errorf("execute reflect template: %w", err)
	}

	return format.Source(buf.Bytes())
}

// Stage generates typescript for all components.
func (tg *TSGen) Stage(dest string, model *Model) error {
	if model == nil {
		return fmt.Errorf("model is nil")
	}

	info, err := os.Stat(dest)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("%q is not a directory", dest)
		}
	} else {
		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(dest, 0777); err != nil {
			return fmt.Errorf("create destination: %w", err)
		}
	}

	sc := stageControl{root: dest}

	if err := sc.writeComponentInterface(tg.Component); err != nil {
		return fmt.Errorf("write component interface: %w", err)
	}
	if err := sc.writeFactoryInterface(tg.ComponentFactory); err != nil {
		return fmt.Errorf("write component factory interface: %w", err)
	}

	for _, c := range model.Components {
		if err := sc.writeComponent(c, func() ([]byte, error) {
			return tg.ComponentConfig(c)
		}); err != nil {
			return fmt.Errorf("write component %s: %w", c.Name, err)
		}
	}

	return nil

}

// Component generates the typescript component interface.
func (tg *TSGen) Component() ([]byte, error) {
	b, err := tsSupportFiles.ReadFile(path.Join(tsSupportDir, "component.ts"))
	if err != nil {
		return nil, err
	}

	return tg.formatTypescript(b)
}

// ComponentFactory generates the typescript component factory interface.
func (tg *TSGen) ComponentFactory() ([]byte, error) {
	b, err := tsSupportFiles.ReadFile(path.Join(tsSupportDir, "component-factory.ts"))
	if err != nil {
		return nil, err
	}

	return tg.formatTypescript(b)
}

// ComponentConfig generates a config interface and factory for a model component.
func (tg *TSGen) ComponentConfig(c Component) ([]byte, error) {
	b, err := tsSupportFiles.ReadFile(path.Join(tsSupportDir, "type-config.ts.tmpl"))
	if err != nil {
		return nil, fmt.Errorf("load type-config template: %w", err)
	}

	t, err := template.New("component").Parse(string(b))
	if err != nil {
		return nil, fmt.Errorf("parse type-config template: %w", err)
	}

	opts := configOptions{
		Name:   c.Name,
		Type:   c.TSName,
		Fields: c.Fields,
		Refs:   c.Referenced(),
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, opts); err != nil {
		return nil, fmt.Errorf("execute type-config template: %w", err)
	}

	return tg.formatTypescript(buf.Bytes())
}

func (tg *TSGen) formatTypescript(in []byte) ([]byte, error) {
	if tg.disableFormat {
		return in, nil
	}

	tf := NewTSFormatter()
	out, err := tf.Format(in)
	if err != nil {
		return nil, fmt.Errorf("unable to format typescript: %w", err)
	}

	return out, nil
}

type reflectOptions struct {
	Names []string
}

type configOptions struct {
	Name   string
	Type   string
	Fields []Field
	Refs   []ImportReference
}

func (co configOptions) RequiredFields() []Field {
	var list []Field
	for _, f := range co.Fields {
		if !f.Optional {
			list = append(list, f)
		}
	}

	return list
}

func (co configOptions) HasOptionalFields() bool {
	return len(co.OptionalFields()) > 0
}

func (co configOptions) OptionalFields() []Field {
	var list []Field
	for _, f := range co.Fields {
		if f.Optional {
			list = append(list, f)
		}
	}

	return list
}

func componentsInPackage(fileSet *token.FileSet, files []*ast.File, importPath string) ([]string, error) {
	p, err := doc.NewFromFiles(fileSet, files, importPath)
	if err != nil {
		return nil, fmt.Errorf("load doc from files: %w", err)
	}

	var list []string

	for _, t := range p.Types {
		if doesCommentContainComponent(t.Doc) {
			list = append(list, t.Name)
		}
	}

	return list, nil
}

var reComponentComment = regexp.MustCompile(`(?m)^\+octant:component$`)

func doesCommentContainComponent(comment string) bool {
	return reComponentComment.MatchString(comment)
}

type stageControl struct {
	root string
}

type stageGenFunc func() ([]byte, error)

func (sc stageControl) writeComponentInterface(fn stageGenFunc) error {
	p := filepath.Join(sc.root, "component.ts")
	return sc.write(fn, p)
}

func (sc stageControl) writeFactoryInterface(fn stageGenFunc) error {
	p := filepath.Join(sc.root, "component-factory.ts")
	return sc.write(fn, p)
}

func (sc stageControl) writeComponent(c Component, fn stageGenFunc) error {
	componentName := strcase.ToKebab(c.TSName)

	p := filepath.Join(sc.root, componentName+".ts")
	return sc.write(fn, p)
}

func (sc stageControl) write(fn stageGenFunc, dest string) error {
	b, err := fn()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dest, b, 0644)
}
