// +build ignore

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	VERSION    = "v0.13.1"
	GOPATH     = os.Getenv("GOPATH")
	GIT_COMMIT = gitCommit()
	BUILD_TIME = time.Now().UTC().Format(time.RFC3339)
	LD_FLAGS   = fmt.Sprintf("-X \"main.buildTime=%s\" -X main.gitCommit=%s", BUILD_TIME, GIT_COMMIT)
	GO_FLAGS   = fmt.Sprintf("-ldflags=%s", LD_FLAGS)
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "build.go",
		Short: "Build tools for Octant",
	}

	rootCmd.Name()

	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "build-electron-dev",
			Short: "electron build",
			Run: func(cmd *cobra.Command, args []string) {
				buildElectronDev()
			},
		},
		&cobra.Command{
			Use:   "ci",
			Short: "full build, running tests",
			Run: func(cmd *cobra.Command, args []string) {
				test()
				vet()
				webDeps()
				webTest()
				webBuild()
				build()
			},
		},
		&cobra.Command{
			Use:   "ci-quick",
			Short: "full build, skipping tests",
			Run: func(cmd *cobra.Command, args []string) {
				webDeps()
				webBuild()
				build()
			},
		},
		&cobra.Command{
			Use:   "web-deps",
			Short: "install client dependencies",
			Run: func(cmd *cobra.Command, args []string) {
				webDeps()
			},
		},
		&cobra.Command{
			Use:   "web-test",
			Short: "run client tests",
			Run: func(cmd *cobra.Command, args []string) {
				webDeps()
				webTest()
			},
		},
		&cobra.Command{
			Use:   "web-build",
			Short: "client build, skipping tests",
			Run: func(cmd *cobra.Command, args []string) {
				webDeps()
				webBuild()
			},
		},
		&cobra.Command{
			Use:   "web",
			Short: "client build, running tests",
			Run: func(cmd *cobra.Command, args []string) {
				webDeps()
				webTest()
				webBuild()
			},
		},
		&cobra.Command{
			Use:   "generate",
			Short: "update generated artifacts",
			Run: func(cmd *cobra.Command, args []string) {
				generate()
			},
		},
		&cobra.Command{
			Use:   "vet",
			Short: "lint server code",
			Run: func(cmd *cobra.Command, args []string) {
				vet()
			},
		},
		&cobra.Command{
			Use:   "test",
			Short: "run server tests",
			Run: func(cmd *cobra.Command, args []string) {
				test()
			},
		},
		&cobra.Command{
			Use:   "build",
			Short: "server build, skipping tests",
			Run: func(cmd *cobra.Command, args []string) {
				build()
			},
		},
		&cobra.Command{
			Use:   "run-dev",
			Short: "run ci produced build",
			Run: func(cmd *cobra.Command, args []string) {
				runDev()
			},
		},
		&cobra.Command{
			Use:   "go-install",
			Short: "install build tools",
			Run: func(cmd *cobra.Command, args []string) {
				goInstall()
			},
		},
		&cobra.Command{
			Use:   "serve",
			Short: "start client and server in development mode",
			Run: func(cmd *cobra.Command, args []string) {
				serve()
			},
		},
		&cobra.Command{
			Use:   "install-test-plugin",
			Short: "build the sample plugin",
			Run: func(cmd *cobra.Command, args []string) {
				installTestPlugin()
			},
		},
		&cobra.Command{
			Use:   "version",
			Short: "",
			Run: func(cmd *cobra.Command, args []string) {
				version()
			},
		},
		&cobra.Command{
			Use:   "release",
			Short: "tag and push a release",
			Run: func(cmd *cobra.Command, args []string) {
				release()
			},
		},
		&cobra.Command{
			Use:   "docker",
			Short: "build octant for use inside a linux container",
			Run: func(cmd *cobra.Command, args []string) {
				docker()
			},
		},
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runCmd(command string, env map[string]string, args ...string) {
	cmd := newCmd(command, env, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	log.Printf("Running: %s\n", cmd.String())
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func runCmdIn(dir, command string, env map[string]string, args ...string) {
	cmd := newCmd(command, env, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	log.Printf("Running in %s: %s\n", dir, cmd.String())
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func newCmd(command string, env map[string]string, args ...string) *exec.Cmd {
	realCommand, err := exec.LookPath(command)
	if err != nil {
		log.Fatalf("unable to find command '%s'", command)
	}

	cmd := exec.Command(realCommand, args...)
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

func buildElectronDev() {
	listCmd := newCmd("go", nil,
		"list", "-m", "-f", "{{.Dir}}", "github.com/asticode/go-astilectron-bundler")

	abPath, err := listCmd.Output()
	if err != nil {
		log.Fatalf("unable to find astilectron-bundler: %w", err)
	}
	abPath = bytes.TrimSpace(abPath)

	runCmdIn(
		filepath.Join("cmd", "octant-electron"),
		"go",
		nil,
		"run",
		filepath.Join(string(abPath), "astilectron-bundler"),
	)
}

func goInstall() {
	pkgs := []string{
		"github.com/GeertJohan/go.rice",
		"github.com/GeertJohan/go.rice/rice",
		"github.com/golang/mock/gomock",
		"github.com/golang/mock/mockgen",
		"github.com/golang/protobuf/protoc-gen-go",
	}
	for _, pkg := range pkgs {
		runCmd("go", map[string]string{"GO111MODULE": "on"}, "install", pkg)
	}
}

func generate() {
	removeFakes()
	runCmd("go", nil, "generate", "-v", "./pkg/...", "./internal/...")
}

func build() {
	newPath := filepath.Join(".", "build")
	os.MkdirAll(newPath, 0755)

	artifact := "octant"
	if runtime.GOOS == "windows" {
		artifact = "octant.exe"
	}
	runCmd("go", nil, "build", "-o", "build/"+artifact, GO_FLAGS, "-v", "./cmd/octant")
}

func runDev() {
	env := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		env[parts[0]] = parts[1]
	}
	runCmd("build/octant", env)
}

func test() {
	runCmd("go", nil, "test", "-v", "./internal/...", "./pkg/...")
}

func vet() {
	runCmd("go", nil, "vet", "./internal/...", "./pkg/...")
}

func webDeps() {
	cmd := newCmd("npm", nil, "ci")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = "./web"
	if err := cmd.Run(); err != nil {
		log.Fatalf("web-deps: %s", err)
	}
}

func webTest() {
	cmd := newCmd("npm", nil, "run", "test:headless")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = "./web"
	if err := cmd.Run(); err != nil {
		log.Fatalf("web-test: %s", err)
	}
}

func webBuild() {
	cmd := newCmd("npm", nil, "run", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = "./web"
	if err := cmd.Run(); err != nil {
		log.Fatalf("web-build: %s", err)
	}
	runCmd("go", nil, "generate", "./web")
}

func serve() {
	var wg sync.WaitGroup

	uiVars := map[string]string{"API_BASE": "http://localhost:7777"}
	uiCmd := newCmd("npm", uiVars, "run", "start")
	uiCmd.Stdout = os.Stdout
	uiCmd.Stderr = os.Stderr
	uiCmd.Stdin = os.Stdin
	uiCmd.Dir = "./web"
	if err := uiCmd.Start(); err != nil {
		log.Fatalf("uiCmd: start: %s", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := uiCmd.Wait(); err != nil {
			log.Fatalf("serve: npm run: %s", err)
		}
	}()

	serverVars := map[string]string{
		"OCTANT_DISABLE_OPEN_BROWSER": "true",
		"OCTANT_LISTENER_ADDR":        "localhost:7777",
		"OCTANT_PROXY_FRONTEND":       "http://localhost:4200",
	}
	serverCmd := newCmd("go", serverVars, "run", "./cmd/octant/main.go")
	serverCmd.Stdout = os.Stdout
	if err := serverCmd.Start(); err != nil {
		log.Fatalf("serveCmd: start: %s", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := serverCmd.Wait(); err != nil {
			log.Fatalf("serve: go run: %s", err)
		}
	}()

	wg.Wait()
}

func installTestPlugin() {
	dir := pluginDir()
	log.Printf("Plugin path: %s", dir)
	os.MkdirAll(dir, 0755)
	filename := "octant-sample-plugin"
	if runtime.GOOS == "windows" {
		filename = "octant-sample-plugin.exe"
	}
	pluginFile := filepath.Join(dir, filename)
	runCmd("go", nil, "build", "-o", pluginFile, "github.com/vmware-tanzu/octant/cmd/octant-sample-plugin")
}

func version() {
	fmt.Println(VERSION)
}

func release() {
	runCmd("git", nil, "tag", "-a", VERSION, "-m", fmt.Sprintf("\"Release %s\"", VERSION))
	runCmd("git", nil, "push", "--follow-tags")
}

func docker() {
	dockerVars := map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        "linux",
		"GOARCH":      "amd64",
	}
	runCmd("go", dockerVars, "build", "-o", "octant", GO_FLAGS, "-v", "./cmd/octant/main.go")
}

func removeFakes() {
	checkDirs := []string{"pkg", "internal"}
	fakePaths := []string{}

	for _, dir := range checkDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				return nil
			}
			if info.Name() == "fake" {
				fakePaths = append(fakePaths, filepath.Join(path, info.Name()))
			}
			return nil
		})
		if err != nil {
			log.Fatalf("generate (%s): %s", dir, err)
		}
	}

	log.Print("Removing fakes from pkg/ and internal/")
	for _, p := range fakePaths {
		os.RemoveAll(p)
	}
}

func gitCommit() string {
	cmd := newCmd("git", nil, "rev-parse", "--short", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		log.Printf("gitCommit: %s", err)
		return ""
	}
	return fmt.Sprintf("%s", out)
}

func pluginDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "octant", "plugins")
	} else if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("LOCALAPPDATA"), "octant", "plugins")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".config", "octant", "plugins")
	}
}
