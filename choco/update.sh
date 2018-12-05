#!/bin/bash

set -e

# This script requires a Github token with repo scope
# It is used to update the chocolately package version

BASE_URL="https://api.github.com/repos/heptio/developer-dash/releases"
GITHUB_URL="https://github.com/heptio/developer-dash/releases"

semver=$(make version)
version=$(echo $semver | cut -c2-)

echo "Latest version is ${semver}" >&2

# Downloading release archives from private GH repo requires an asset ID
# Get asset IDs from archive
name="hcli_${version}_Windows-64bit.zip"
response="$BASE_URL/tags/${semver}?access_token=$GITHUB_TOKEN"
eval $(curl -L $response | grep -C3 "name.:.\+$name" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')
[ "$id" ] || { echo "Error: Failed to get asset id, response: $response" | awk 'length($0)<100' >&2; exit 1; }

download="$BASE_URL/assets/$id?access_token=$GITHUB_TOKEN"
checksum64=$(curl --fail -L -H 'Accept: application/octet-stream' "${download}" | shasum -a 256 - | cut -f 1 -d " ")

url64="$GITHUB_URL/download/${semver}/$name"

# If public, use this instead
#url64="https://github.com/heptio/developer-dash/releases/download/$(make version)/hcli_$(make version | cut -c2-)_Windows-64bit.zip"
#
#checksum64=$(curl --fail -L "${url64}" | shasum -a 256 - | cut -f 1 -d " ")

echo "Updating choco metadata..." >&2
sed -i.bak "s/<version>.*<\/version>/<version>${version}<\/version>/" ./choco/hcli.nuspec
sed -i.bak "s!^\$url64 = '.*'!\$url64 = '${url64}'!" ./choco/tools/chocolateyinstall.ps1
sed -i.bak "s/^\$checksum64 = '.*'/\$checksum64 = '${checksum64}'/" ./choco/tools/chocolateyinstall.ps1

