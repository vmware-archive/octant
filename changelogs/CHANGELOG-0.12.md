## v0.12.0
#### 2020-04-22

### Download
 - https://github.com/vmware/octant/releases/tag/v0.12.0

### Highlights
  * Add new vertical navigation (@mklanjsek, #835)
  * Add breadcrumbs, reworked headers for all pages (#710, @mklanjsek)
  * Add support for Clarity's single action for data grid rows (#801, @bryanl)
  * Add support for creating objects (#802, @bryanl)
  * Update logging interfaces to support combined loggint and streaming content. (#637, @wwitzel3)
  * Add support for different persistent volume sources in summary view (#817, @GuessWhoSamFoo)
  * Add network policy printer (#813, @GuessWhoSamFoo)

### Breaking API Changes
  * **PLUGIN API** Add ActionName to service.ActionRequest to allow dispatching of Octant events (#808, @wwitzel3)
  * Add namespacing to Octant actions and events (#820, @wwitzel3)

### All Changes
  * Refactor backend tab creation. (#861, @bryanl)
  * Improve UI link rendering time. (#854, @bryanl)
  * Fix issues with custom resource breadcrumbs (#855, @mklanjsek)
  * Fix truncation issues with namespace and context selectors (#737, @mklanjsek)
  * Fix logs console errors when timestamp are not present (#840, @mklanjsek)
  * Fix nightly builds that broke when we moved to GitHub Actions (#839, @wwitzel3)
  * Fix bug that prevented Pod summary loading if status was nil/unknown. (#831, @wwitzel3)
  * Fix issues with custom resource breadcrumbs
  * Enhance multiple change detections with different values from pipe
  * Decrease link latency
  * Fix truncation issues with namespace and context selectors
  * Fix logs console errors when no timestamp
  * Update ngx-highlightjx to ivy compatible version
  * Remove setup-go path steps and bump setup-node versions
  * Update karma and protractor versions
  * Add unit tests for navigation service
  * Update nightly with tag and tmp file
  * Load summary for pods with nil status
  * Improve selection indication for vertical nav
  * Add ellipsis to labels when they are long to avoid overflow
  * Fix labels/selectors overflow their containers
  * Add plugin topic link to README.md
  * Enhance heading anchors actually visible/usable
  * Fix error logging causing console spam
  * Add describer for generalized network policy cases
  * Fix event types for remaining handlers
  * Fix datagrid filtering
  * Fix frontend handler assignment
  * Add more context to Overview links
  * Support different volume sources in summary view
  * Fix go mod issues by ignoring go build and doc main files
  * Remove old reference to CLA
  * Add e2e test proposal
  * Run prettier on angular app dir
  * Extract creation of frontend/backend handlers
  * Update client side dependencies
