# Chocolatey Packages

This directory contains chocolatey packages that can be pushed to the public registry in the future. CI steps will build a NuGet package from the an existing `*.nuspec` file then install locally using chocolatey for testing.

Eventually this should be move to a separate repository to track multiple chocolately packages.

## Requirements

chocolatey0.10.11

## 1. Scaffold a new package

This generates a new directory containing a `*.nuspec` file and PowerShell scripts for install/uninstall.

```
choco new <package-name>
```

Refer to the documentation for flags to prefill metadata.

## 2. Edit package metadata

The `*.nuspec` file contains information about the package maintainers, project website, licenses, etc. Instructions are commented out within the file and remove them as fields are completed.

Inside the tools directory, there is a PowerShell installer script. Edit the URL to the location of the `*.exe` or `*.zip`.

## 3. Build the package

Generate a NuGet package with the `*.nupkg` extension from the `*.nuspec` file.

```
choco pack <name>.nuspec
```

## 4. Push the package to a chocolatey server

Adding a package to `chocolatey` will require an account with an API key. Chocolatey will also undergo an additional review of the package before it is publically available.

```
choco push <name>.nupkg --apikey="<key>" --source="http://localhost/chocolatey"
```

## (Optional) Add a chocolatey server

```
choco source add --name=<name> --source="http://localhost/chocolatey"
```

## Resources

[Setting up a testing environment with Vagrant](https://github.com/chocolatey/chocolatey-test-environment)
[Workshop and exercises](https://github.com/ferventcoder/chocolatey-workshop)
[Automatic package updates](https://github.com/majkinetor/au)
