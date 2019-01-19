This package provides a way to bundle an [astilectron](https://github.com/asticode/go-astilectron) app using the [bootstrap](https://github.com/asticode/go-astilectron-bootstrap).

Check out the [demo](https://github.com/asticode/go-astilectron-demo) to see a working example.

# Installation

Run the following command:

    $ go get -u github.com/asticode/go-astilectron-bundler/...

# Build the binary
    
Run the following command:

    $ go install github.com/asticode/go-astilectron-bundler/astilectron-bundler
    
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

## Custom paths

You can set the following paths:

- `input_path`: path to your project. defaults to the current directory
- `go_binary_path`: path to the `go` binary. defaults to "go"
- `output_path`: path to the dir where you'll find the bundle results. defaults to `current directory/output`
- `resources_path`: path where the `resources` dir is and will be written. path must be relative to the `input_path`. defaults to "resources"
- `vendor_dir_path`: path where the `vendor` dir will be written. path must be relative to the `output_path`
- `working_directory_path`: path to the dir where the bundler runs its operations such as provisioning the vendor files or binding data to the binary

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

# Usage

If **astilectron-bundler** has been installed properly (and the $GOPATH is in your $PATH), run the following command:

    $ astilectron-bundler -v -c <path to your configuration file>
    
or if your working directory is your project directory and your bundler configuration has the proper name (`bundler.json`)

    $ astilectron-bundler -v
    
# Output

For each environment you specify in your configuration file, **astilectron-bundler** will create a folder `<output path you specified in the configuration file>/<os>-<arch>` that will contain the proper files.

# Ldflags

**astilectron-bundler** uses `ldflags` when building the project. It means if you add one of the following variables as global exported variables in your project, they will have the following value:

- `AppName`:  filled with the configuration app name
- `BuiltAt`: filled with the date the build has been done at

If you need to add more flags yourself, like for a version number, add something
like this to your `astilectron-bundler` command: `-ldflags X:main.Version=xyzzy`.

If you need to add multiple flags you can pass `-ldflags` multiple times, with
multiple values split on commas, like this:

`-ldflags X:main.Version=xyzzy,main.CommitCount=100 -ldflags race`

That would set two variables and enable the race detection.

# Subcommands
## Only bind data: bd

Use this subcommand if you want to skip most of the bundling process and only bind data/generate the `bind.go` file (useful when you want to test your app running `go run *.go`):

    $ astilectron-bundler bd -v -c <path to your configuration file>

## Clear the cache: cc

The **bundler** stores downloaded files in a cache to avoid downloading them over and over again. That cache may be corrupted. In that case, use this subcommand to clear the cache:

    $ astilectron-bundler cc -v
