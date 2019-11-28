# Getting Started

## Environment Variables

Octant is configurable through environment variables defined at runtime, here are some of the notable variables:

### User Variables

* `KUBECONFIG` - set to non-empty location if you want to set KUBECONFIG with an environment variable.
* `OCTANT_NAMESPACE` - initial namespace to load when Octant starts.
* `OCTANT_CONTEXT` - intial context to load when Octant starts.
* `OCTANT_DISABLE_CLUSTER_OVERVIEW` - disable cluster overview when a context does not have cluster level permissions.
* `OCTANT_PLUGIN_PATH` - add a plugin directory or multiple directories separated by `:`. Plugins will load by default from `$HOME/.config/octant/plugins`


**Notice:** If using [fish shell](https://fishshell.com), tilde expansion may not occur when using `env` to set environment variables.

### Flags as Variables

All command-line flags can also be passed as environment variables by using all UPPERCASE, replacing the `-` with `_` and prefixing them with `OCTANT_`.
When using CLI flags that enable/disable a feature the following values are considered true and false:

  * **True** - "1", "t", "T", "true", "TRUE", "True"
  * **False** - "0", "f", "F", "false", "FALSE", "False"

Example:

 * `--namespace=default` becomes `OCTANT_NAMESPACE=default`

## Command Line Flags

Octant is configurable through command line flags set at runtime. You can see all of the available options by
running `octant --help`.

### Verbosity

The verbosity has a special type that is used to parse the flag, which means it can be provided
shorthand by just adding more `v` to equal the level count or with an explicit equal sign.

```sh
-v[vv], --verbosity=count      verbosity level
```

For example

```sh
octant -vvv
```

Is equal to

```sh
octant --verbosity=3
```