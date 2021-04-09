## v0.19.0
#### 2021-04-08

### Download
 - https://github.com/vmware-tanzu/octant/releases/v0.19.0

### Breaking API Changes
  * Moved `request.ClientID` to `request.ClientState` (#2244, @xtreme-vikram-yadav)
  * Changed button groups to accept button components (@2255, @ftovaro)

### Highlights
  * Upgraded project to Clarity 5 (#2222, @mklanjsek)
  * Shared partial octant state with plugins (#2244, @xtreme-vikram-yadav)
  * Created ButtonComponent and extended LinkComponent to receive components to be wrapped as links (#2255, @ftovaro)
  * Added Delete and Create object store calls for Go Plugins (#2257, @xtreme-vikram-yadav)

### Bug Fixes
  * Fixed expression selector rendering (#2252, @xtreme-jon-ji)
  * Fixed changelog generation script to support usernames containing '-' (#2270, @xtreme-jon-ji)
  * Fixed Resource Viewer bug with missing pods (#2280, @mklanjsek)

### All Changes
  * Refactored forms to improve re-usability (#1504, @lenriquez)
  * Added Signpost (#2018, @lenriquez)
  * Added ability to change dynamic component on the fly (#2215, @mklanjsek)
  * Added CreateLink method with objects for go plugin (#2276, @GuessWhoSamFoo)

