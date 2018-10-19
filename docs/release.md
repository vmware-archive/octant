## Steps to release hcli via TravisCI

`GITHUB_TOKEN` needs to be set from the TravisCI UI with a Github token of the proper scope (currently private repo but eventually should be public only).

1. The version is tracked in the Makefile. Update and merge the new version prior to release.

2. Pull the change then run `make release`

3. A new build will be triggered then a draft of the artifacts will be available for review in Github Releases.

## Steps to release hcli manually

What you'll need:

 - a Github token with repo scope enabled
 - an installation of rpmbuild
 - an installation of [goreleaser](https://goreleaser.com)

Start the release process with the commands:

```
VERSION=v0.0.1
git tag -a $VERSION -m "$VERSION release"
git push origin $VERSION
goreleaser --rm-dist
```

This will
 - create a [semver](http://semver.org) git tag
 - tag the latest commit
 - push the tag to the repo
 - build and push the release
