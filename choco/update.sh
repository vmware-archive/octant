#!/bin/bash

set -e

BASE_URL="https://api.github.com/repos/vmware-tanzu/octant/releases"
GITHUB_URL="https://github.com/vmware-tanzu/octant/releases"

semver=$(go run build.go version)
version=$(echo $semver | cut -c2-)

echo "Latest version is ${semver}" >&2

# Downloading release archives from private GH repo requires an asset ID
# Get asset IDs from archive
name="octant_${version}_Windows-64bit.zip"

url64="https://github.com/vmware-tanzu/octant/releases/download/${semver}/octant_${version}_Windows-64bit.zip"

checksum64=$(curl --fail -L "${url64}" | shasum -a 256 - | cut -f 1 -d " ")

echo "Updating choco metadata..." >&2
sed -i.bak "s/<version>.*<\/version>/<version>${version}<\/version>/" ./choco/octant.nuspec
sed -i.bak "s!^\$url64 = '.*'!\$url64 = '${url64}'!" ./choco/tools/chocolateyinstall.ps1
sed -i.bak "s/^\$checksum64 = '.*'/\$checksum64 = '${checksum64}'/" ./choco/tools/chocolateyinstall.ps1

