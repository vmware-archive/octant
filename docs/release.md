## Steps to release hcli

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
