- [v0.11.0](#v0110)

## v0.11.0
#### 2020-03-20

### Download
 - https://github.com/vmware/octant/releases/tag/v0.11.0

### Highlights
  * Upgraded Octant to Angular9 / Clarity3 (#567, @wwitzel3)
  * Created shared module for angular frontend (#607, @bryanl)
  * Created Electron application stub (#619, @bryanl)
  * Created persistent volume printer (#684, @dotNomad)
  * Added Log search / filter (#713, @mklanjsek)
  * Replaced existing YAML viewer with read-only editor using Monaco (#723, @wwitzel3)
  * Moved internal `Logger` to `pkg/log` (#761, @dotNomad)
  * Moved terminal to tab (#775, @mklanjsek)

### Breaking API Changes
  * Removed returned boolean value from `store.Get` (#643, @GuessWhoSamFoo)
  * Changed `service.Request` to an interface and `Path` property accessed via `Path()` (#645, @wwitzel3)

### All Changes
  * Fixed issue where plugins would not load on Windows. (#580, @wwitzel3)
  * Added new component for multi-line text to handle long ConfigMap values (#583, @GuessWhoSamFoo)
  * Allow for web and app UIs to work simultaneously (#610, @bryanl)
  * Fixed annotation and log components to be accessible by plugins (#618, @GuessWhoSamFoo)
  * Remove extra actions from poller and websocket client (#631, @bryanl)
  * Added `ListNamespaces` for plugins (#630, @GuessWhoSamFoo)
  * Changed labels and selectors to use clarity labels (#644, @GuessWhoSamFoo)
  * Fixed container logs when the log content does not contain a timestamp. (#649, @wwitzel3)
  * Fixed the log viewer size (#666, @mklanjsek)
  * Changed Log view to scroll to bottom at initial load (#667, @mklanjsek)
  * Removed white borders from icons (#678, @bryanl)
  * Added cordon and uncordon actions (#680, @GuessWhoSamFoo)
  * Fixed issue where logs stopped working after context-switching  (#683, @mothershipper)
  * Do not show empty CRD versions (#688, @bryanl)
  * Fixed delay when showing datagrids with filters (#692, @bryanl)
  * Don't calculate pagination for data grid until there are rows (#693, @bryanl)
  * Widened Log timestamp so it doesn't break up to the second line (#700, @mklanjsek)
  * Added log filtering (including Regex) and support case sensitive search (#702, @mklanjsek)
  * Fixed bug with Terminal exec not working on Azure Kubernetes Service Virtual Nodes (#707, @wwitzel3)
  * Configured grpc client to wait for server to be ready before sending requests (#713, @bryanl)
  * Added status indicators to be shown by text components (#715, @bryanl)
  * Changed summary sections can be updated (#717, @bryanl)
  * Added support for ANSI escape codes in log viewer (#722, @mklanjsek)
  * Added manual cronjob trigger action (#731, @GuessWhoSamFoo)
  * Moved delete button next to page header (#749, @GuessWhoSamFoo)
  * Fixed bug where CRD links stop working after switching contexts (#751, @GuessWhoSamFoo)
  * Added support for globs in ingress TLS host (#755, @bryanl)
  * Added CSS to change namespace selector colors when active (#780, @GuessWhoSamFoo)
