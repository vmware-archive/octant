 - [v0.6.0](#v060)

## v0.6.0
#### 2019-08-20

### Download
 - https://github.com/vmware-tanzu/octant/releases/tag/v0.6.0

### Highlights
- Fixed cases where CRDs were causing errors from an int64 to float64 conversion
- Octant now starts on port 7777 instead of a random port
- Loading indicators are available to give users better feedback when loading a list of resources
- Many improvements and bug fixes since initial release

### Bug Fixes
  * Fixed cluster or namespaced scoped CRDs sometimes not showing up (#146, @bryanl)
  * Fixed sample installation Makefile target (#151, @alexmt)
  * Fixed bug causing octant to 403 when setting `OCTANT_LISTENER_ADDR` (#152, @wwitzel3)
  * Fixed scroll getting stuck when viewing container logs (#162, @nfarruggiagl)
  * Fixed initial set namespace flag by respecting initial URL routing (#165, @GuessWhoSamFoo)
  * Fixed bug where switching clusters with CRDs loaded causes octant to crash (#170, @bryanl)
  * Fixed displaying external IP or hostname when service is exposed (#197, @GuessWhoSamFoo)

### Enhancements
  * Added links from pods to its node (#115, @bryanl)
  * Added int64 to float64 conversation for unstructured converter (#145, @bryanl)
  * Added Init container labels to deployments to match pods (#148, @GuessWhoSamFoo)
  * Added support to multiple kubeconfig files (#164, @bryanl)
  * Added loading indicators to objects in nav when their informer has not synced (#176, @bryanl)
  * Changed the default listening port to `7777` (#185, @wwitzel3)
  * Added pods with succeeded phase to be listed under job viewer (#190, @GuessWhoSamFoo)
  * Added environment variable to set accepted hosts (#194, @aksalj)

### Documentation
  * Fixed broken documentation links (#128, @schallert)
  * Clarified readme installation for various operating systems (#133, @nickgerace)
  * Added note for users to expect octant to launch in a browser on a given port (#138, @mikeroySoft)
  * Clarified running octant on a given host and port (#157, @sensay-nelson)
