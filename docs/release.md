## Steps to release sugarloaf via TravisCI (Deprecated)

`GITHUB_TOKEN` needs to be set with the "repo" scope to have access to private repositories. This token should be regenerated once the project is public. See [TravisCI docs](https://docs.travis-ci.com/user/travis-ci-for-private/#how-can-i-make-a-private-repository-public) for changing from a private to public project.

1. The version is tracked in the Makefile. Update and merge the new version prior to release.

2. Pull the change then run `make release`

3. A new build will be triggered then a draft of the artifacts will be available for review in Github Releases.

### Encrypt a token

1. Download the Travis CLI ruby gem: `gem install travis`

2. Login with your Github account: `travis login --pro`.

3. Generate the token through the heptibot Github account. Navigate to the project root and encrypt the token: `travis encrypt GITHUB_TOKEN="..." --add`. The token can only be decrypted by Travis CI, not the encrypter or owners of the repository.

4. Commit the changes to `.travis.yml` and push.

## Steps to release sugarloaf manually

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
