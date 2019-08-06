 - [v0.5.1](#v051)
 - [v0.5.0](#v050)

## v0.5.1
#### 2019-08-06

### Download
 - https://github.com/vmware/octant/releases/tag/v0.5.1

### Bug fixes
 - Fixed bug with context switcher stream registration (#110, @wwitzel3)

## v0.5.0
#### 2019-08-05

### Download
 - https://github.com/vmware/octant/releases/tag/v0.5.0

### Highlights
- Significantly reduced memory usage and faster performance
- Summary and YAML view for nodes are now available under Cluster Overview

### All Changes:

  * Reworked deployment status so it is consistent with other objects (#92, @bryanl)
  * Removed the Watch store to help improve performance and memory usage (#94, @wwitzel3)
  * Added pagination to data grids (#103, @bryanl)
  * Increased REST client's QPS (#102, @bryanl)
  * Changed the store API List function to return an unstructured list (#102, @bryanl)
  * Created an `allow` list to reduce the number of objects that need to be checked when building the resource viewer (#102, @bryanl)
  * Used unstructured objects in more places to removes a few unnecessary conversions (#102, @bryanl)
  * Added nodes to cluster overview (#104, @bryanl)
  * Changed `--kubeConfig` flag to `--kubeconfig` (#107, @GuessWhoSamFoo)
