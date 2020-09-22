This package provides a way to bundle an [astilectron](https://github.com/asticode/go-astilectron) app using the [bootstrap](https://github.com/asticode/go-astilectron-bootstrap).

Check out the [demo](https://github.com/asticode/go-astilectron-demo) to see a working example.

# Installation

Run the following command:

```shell
go get -u github.com/asticode/go-astilectron-bundler/...
```

# Build the binary

Run the following command:

```shell
go install github.com/asticode/go-astilectron-bundler/astilectron-bundler
```

# Configuration

**astilectron-bundler** uses a configuration file to know what it's supposed to do.

## Basic configuration

Here's the basic configuration you'll usually need:

```json
{
  "app_name": "Test",
  "icon_path_darwin": "path/to/icon.icns",
  "icon_path_linux": "path/to/icon.png",
  "icon_path_windows": "path/to/icon.ico"
}
```

It will process the project located in the current directory and bundle it in the `output` dir for your os/arch.

## Bundle for specific Astilectron and/or Electron versions

The following customization can be made to `bundler.json`

* `version_electron` - version of electron, defaults to the value specified in the `go-astilectron` version you're using
* `version_astilectron` - version of astilectron, defaults to the value specified in the `go-astilectron` version you're using

## Bundle for other environments

You can bundle your project for multiple environments with the `environments` key:

```json
{
  "environments": [
    {"arch": "amd64", "os": "darwin"},
    {"arch": "amd64", "os": "linux"},
    {
      "arch": "amd64",
      "os": "windows",
      "env": {
        "CC": "x86_64-w64-mingw32-gcc",
        "CXX": "x86_64-w64-mingw32-g++",
        "CGO_ENABLED": "1"
      }
    }
  ]
}
```

For each environment you can specify environment variables with the `env` key.

## Adapt resources

You can execute custom actions on your resources before binding them to the binary such as uglifying the `.js` files with the `resources_adapters` key:

```json
{
  "resources_adapters": [
    {
      "args": ["myfile.js", "mynewfile.js"],
      "name": "mv"
    },
    {
      "args": ["-flag", "value", "mynewfile.js"],
      "name": "myawesomebinary"
    }
  ]
}
```

All paths must be relative to the `resources` folder except if you provide a `dir` option (a path relative to the `resources` folder) in which case it will be relative to that path.

## Build flags

You can pass arbitrary build flags into the build command with the `build_flags` key:

```json
{
  "build_flags": {
    "gcflags": "\"all=-N -l\""
  }
}
```

## Custom paths

You can set the following paths:

* `input_path`: path to your project. defaults to the current directory
* `go_binary_path`: path to the `go` binary. defaults to "go"
* `output_path`: path to the dir where you'll find the bundle results. defaults to `current directory/output`
* `resources_path`: path where the `resources` dir is and will be written. path must be relative to the `input_path`. defaults to "resources"
* `vendor_dir_path`: path where the `vendor` dir will be written. path must be relative to the `output_path`
* `working_directory_path`: path to the dir where the bundler runs its operations such as provisioning the vendor files or binding data to the binary

## Adapt the bind configuration

You can use the `bind` attribute to alter the bind configuration like so:

```json
{
  "bind": {
    "output_path": "path/to/bind/output/path",
    "package": "mypkg"
  }
}
```

* `output_path`: path to the directory where you want bind files to be created. defaults to the current working directory
* `package`: the package name to use for the bind files. defaults to "main"

When you specify an `output_path`, the `package` will **probably** need to be set.

## Info.plist generation from the bundler configuration file property

You can add custom **Info.plist** configuration to the **bundler.json**:

```json
{
  "app_name": "Best App",
  "icon_path_darwin": "resources/icon.icns",
  "info_plist": {
    "CFBundlePackageType": "APPL",
    "CFBundleInfoDictionaryVersion": "6.0",
    "CFBundleIconFile": "icon.icns",
    "CFBundleDisplayName": "Best App",
    "CFBundleExecutable": "app_binary",
    "CFBundleIdentifier": "com.company.BestApp",
    "LSUIElement": "NO",
    "LSMinimumSystemVersion": "10.11",
    "NSHighResolutionCapable": true,
    "NSAppTransportSecurity": {
      "NSAllowsArbitraryLoads": true
    }
  }
}
```

# Usage

If **astilectron-bundler** has been installed properly (and the $GOPATH is in your $PATH), run the following command:

```shell
astilectron-bundler -c <path to your configuration file>
```

or if your working directory is your project directory and your bundler configuration has the proper name (`bundler.json`)

```shell
astilectron-bundler
```

## Output

For each environment you specify in your configuration file, **astilectron-bundler** will create a folder `<output_path you specified in the configuration file>/<os>-<arch>` that will contain the proper files.

# Ldflags

**astilectron-bundler** uses `ldflags` when building the project. It means if you add one of the following variables as global exported variables in your project, they will have the following value:

* `AppName`:  filled with the configuration app name
* `BuiltAt`: filled with the date the build has been done at
* `VersionAstilectron`: filled with the version of Astilectron being bundled/used
* `VersionElectron`: filled the version of Electron being bundled/used

You can use the following to alter the Ldflags behavior:

```json
{
  "ldflags_package": "some/path/to/pkg"
}
```

* `ldflags_package`: which local package these variables exist in. defaults to `bind`'s `package` value (for backwards compatibility)

If you need to add more flags yourself, like for a version number, add something
like this to your `astilectron-bundler` command: `-ldflags X:main.Version=xyzzy`.

If you need to add multiple flags you can pass `-ldflags` multiple times, with
multiple values split on commas, like this:

```shell
-ldflags X:main.Version=xyzzy,main.CommitCount=100 -ldflags race
```

That would set two variables and enable the race detection.

(in either case, make sure to substitute `main` with the package where your `Version`/`CommitCount`/etc. variables exist)

# Commands

## Only bind data: bd

Use this command if you want to skip most of the bundling process and only bind data/generate the `bind.go` file (useful when you want to test your app running `go run *.go`):

```shell
astilectron-bundler bd -c <path to your configuration file>
```

## Clear the cache: cc

The **bundler** stores downloaded files in a cache to avoid downloading them over and over again. That cache may be corrupted. In that case, use this command to clear the cache:

```shell
astilectron-bundler cc
```

# Frequent problems

## "xxx architecture of input file `xxx' is incompatible with xxx output"

When building for `linux` you may face an error looking like this:

```shell
FATA[0009] bundling failed: bundling for environment linux/amd64 failed: building failed: # github.com/asticode/go-astilectron-demo
/usr/local/go/pkg/tool/linux_amd64/link: running gcc failed: exit status 1
/usr/bin/ld: i386 architecture of input file `/tmp/go-link-275377070/000000.o' is incompatible with i386:x86-64 output
collect2: error: ld returned 1 exit status
```

Thanks to [this comment](https://github.com/asticode/go-astilectron-demo/issues/28#issuecomment-509050603), you need to add the `ldflags` key to your `bundler.json` with the value `{"linkmode":["internal"]}`.
