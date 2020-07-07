// +build ignore

/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
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
	flag.Parse()
	for _, cmd := range flag.Args() {
		switch cmd {
		case "build-electron-dev":
			buildElectronDev()
		case "ci":
			test()
			vet()
			webDeps()
			webTest()
			webBuild()
			build()
		case "ci-quick":
			webDeps()
			webBuild()
			build()
		case "web-deps":
			webDeps()
		case "web-test":
			webDeps()
			webTest()
		case "web-build":
			webDeps()
			webBuild()
		case "web":
			webDeps()
			webTest()
			webBuild()
		case "clean":
			clean()
		case "generate":
			generate()
		case "vet":
			vet()
		case "test":
			test()
		case "build":
			build()
		case "run-dev":
			runDev()
		case "go-install":
			goInstall()
		case "serve":
			serve()
		case "install-test-plugin":
			installTestPlugin()
		case "version":
			version()
		case "release":
			release()
		case "docker":
			docker()
		default:
			log.Fatalf("Unknown command %q", cmd)
		}
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

func clean() {
	if err := os.Remove("pkg/icon/rice-box.go"); err != nil {
		log.Fatalf("clean: %s", err)
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
