## v0.16.2
#### 2020-11-11

### Download
 - https://github.com/vmware-tanzu/octant/releases/v0.16.2
 
### All Changes
  * Log warning when external metrics unavailable (#1023, @wwitzel3)
  * Inform user that the backend uses an annotation when AWS ALB use-annotation is set (#1060, @wwitzel3)
  * Fixed panic when websocket ClientFromRequest returns an error (#1340, @wwitzel3)
  * Use link when displaying Ingress host (#1431, @mklanjsek)
  * Removed formGroup from payload actions and updated to return values for multiple checkboxes (#1563, @GuessWhoSamFoo)
  * Prevent angular sanitizer from removing log content (#1582, @mklanjsek)
  * Show ConfigMaps and PersistentVolumeClaims in Pod's resource viewer (#1541 @alexbrand)
  * Fixed the first/last seen to display null instead of a pointer (#1386 @sladyn98)
  * Switched resolved address to pubic npm registry (#1571 @scothis)
  * Moved websocket service to new module (#1459 @bryanl)
  * Implemented Help menu additional items (#875 @mklanjsek)

## v0.16.1
#### 2020-10-07

### Download
 - https://github.com/vmware-tanzu/octant/releases/v0.16.1

### All Changes
  * Fixed resource viewer unmarshal bug (#1424, @scothis)
  * Added a streamer interface for logger (#1425, @bryanl)
  * Fixed syntax highlighting in YAML tab (#1436, @GuessWhoSamFoo)
  * Changed default internal DashService API to 127.0.0.1 rather than localhost (#1438, @joshrosso)
  * Fixed clientID not added to ActionRequest (#1441, @GuessWhoSamFoo)
  * Moved breadcrumb/title generation to tab component (#1445, @GuessWhoSamFoo)
  * Fixed regression with  multiple plugin tabs overwriting (#1446, @GuessWhoSamFoo)
  * Added `trustedContent` flag for Markdown (#1460, @wwitzel3)
  * Changed redirect path when switching namespaces to be more intuitive (#1461, @GuessWhoSamFoo)
  * Fixed logs not scrolling to bottom when selecting container or time (#1464, @GuessWhoSamFoo)

## v0.16.0
#### 2020-09-24

### Download
 - https://github.com/vmware-tanzu/octant/releases/v0.16.0

### Highlights
  * Added a 404-style error page when a resource is not found (#422, @scothis)
  * Changed default log viewer to show last 5 minutes and allow selecting a broader range (#1209, @wwitzel3)
  * Added `SendAlert` to plugin interface (#1216, @GuessWhoSamFoo)
  * Changed to dynamic component loading (#1242, @bryanl)
  * Added `Ctrl+/` keyboard shortcut and list of keyboard shortcuts (#1319, @wwitzel3)
  * Updated quick switcher UI and added namespace to search (#1381, @GuessWhoSamFoo)

### Bug Fixes
  * Fixed default provided namespaces to initial namespace when empty (#838, @wwitzel3)
  * Fixed problem with Storybook rendering of dynamic components (#1289, @mklanjsek)
  * Fixed editing service to show sorted selectors (#1302, @GuessWhoSamFoo)
  * Fixed safari height bug in the header so all browsers render the header the same. (#1313, @alexbarbato)
  * Fixed configuring GRPC message size to API Client (#1324, @nodece)
  * Fixed compiler warnings by colidating SCSS dependencies (#1357, @mklanjsek)
  * Fixed pv list generation when claimRef pvc cannot be found (#1358, @GuessWhoSamFoo)
  * Fixed unreferenced ConfigMap crashing summary tab (#1362, @GuessWhoSamFoo)

### All Changes
  * Exposed full selector capabilities through `Key` object (#1201, @ipsi)
  * Updated build to use Golang 1.15 (#1248, @scothis)
  * Added `SendEvent` support to JavaScript plugin runtime (#1290, @wwitzel3)
  * Added an optional button group to data grid tables (#1299, @scothis)
  * Added validator and action payload for stepper (#1300, @GuessWhoSamFoo)
  * Added support for selectors in JavaScript plugin dashboard client (#1304, @bryanl)
  * Added modal component and opening modals through buttons (#1305, @GuessWhoSamFoo)
  * Added an Octant log sink for zap message (#1321, @bryanl)
  * Changed to icons and colors for indicator to make it more accessible (#1335, @wwitzel3)
