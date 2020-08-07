/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vmware-tanzu/octant/pkg/view/tsgen"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	defaultSourceValue := filepath.Join(dir, "pkg", "view", "component")

	var source string
	flag.StringVar(&source, "source", defaultSourceValue, "source dir")
	var dest string
	flag.StringVar(&dest, "dest", "", "destination dir")
	flag.Parse()

	if err := run(source, dest); err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
}

func run(source, dest string) error {
	tg, err := tsgen.NewTSGen()
	if err != nil {
		return fmt.Errorf("create typescript generator: %w", err)
	}

	names, err := tg.Names(source)
	if err != nil {
		return fmt.Errorf("get component names: %w", err)
	}

	m, err := tg.Reflect(names)
	if err != nil {
		return fmt.Errorf("run reflect: %w", err)
	}

	if err := tg.Stage(dest, m); err != nil {
		return fmt.Errorf("stage typescript: %w", err)
	}

	return nil
}
