## v0.13.0
#### 2020-06-01

### Download
 - https://github.com/vmware-tanzu/octant/releases/tag/v0.13.0

### Highlights
  * Added resource editing
  * Added datagrid actions
  * Added resource deleting indicators
  * Added storybook for UI prototying and exploring
  * Added Electron prototype
  * Upgraded client-go and remove k8s.io/kubernetes imports
  * Upgraded to Angular 9.1.9 and Clarity 3.1.3

### All Changes
  * Fixed Terminal window resizing issues (#935, @mklanjsek)
  * Fixed Graphviz component performance (#951, @mklanjsek)
  * Fixed problems with the graphviz component (#948, @mklanjsek)
  * Added informer error handler using latest client-go (#942, @wwitzel3)
  * Added `--namespace-list` flag to accept a slice of namespaces (#928, @GuessWhoSamFoo)
  * Added terminal support for multiple containers (#916, @GuessWhoSamFoo)
  * Added tooltips for selectors and labels (#884, @mklanjsek)
  * Updated terminal manager and map terminal instances to each websocket client (#887, @GuessWhoSamFoo)
  * Fixed case where unclaimed persistent volumes causes panic (#889, @GuessWhoSamFoo)
  * Added missing label selector to protobuf key request (#908, @GuessWhoSamFoo)
  * Fixed replica spec to checked for desired replica count
  * Added helper menu with build information
  * Fixed Terminal window resizing issues
  * Added Datagrid actions
  * Added Resource deleting indicators
  * Added Storybook for UI prototying and exploring
  * Added Electron prototype
  * Upgraded client-go and remove k8s.io/kubernetes imports
  * Upgraded to Angular 9.1.9 and Clarity 3.1.3
  * Fixed graphviz component
  * Fixed dark mode colors for log viewer
  * Removed duplicated titles and color from visited breadcrumb links
  * Fixed missing context from namespace get
  * Updated schema.Convert tests and data
  * Fixed terminal exit corner cases
  * Upgrade gomock to 1.4.3
  * Removed k8s.io/kubernetes from deps
  * Fixed potential race in config watcher test
  * Removed old style icons and icon package
  * Added CLI argument `browser-path`
  * Fixed problems with timestamp tooltips and overflow labels/selectors
  * Added delete action to object lists
  * Fixed client caching, no-cache, no-store headers for root (index) request
  * Fixed missing label selector to protobuf key request
  * Added Show service account secrets
  * Fixed a couple of Tooltips issues on Firefox
  * Fixed case where unclaimed persistent volume causes panic
  * Fixed container log tests 
  * Added editing of objects
  * Fixed layout MarkdownText better in tables/cards
  * Added logging UI improvements
  * Added ellipsis to selector overflow pattern & fix lint issues
