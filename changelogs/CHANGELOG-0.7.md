- [v0.7.0](#v070)

## v0.7.0
#### 2019-09-18

### Download
 - https://github.com/vmware/octant/releases/tag/v0.7.0

### Highlights
 - Replaced event streams with websockets (#239, @bryanl)
 - Added support for deleting pods (#227, @bryanl)

### Bug Fixes
  * Fixed bug where mounted pods of persistent volume claims were not displayed (#281, @GuessWhoSamFoo)
  * Removed informer for custom resources when they no longer exist (#250, @bryanl)
  * Changed Deployment printer to use ReplicaSet owner reference to get pods (#245, @GuessWhoSamFoo)
  * Fixed bug where tabs generated from plugins would switch back to summary (#240, @GuessWhoSamFoo)
  * CI now runs TS linting (#224, @wwitzel3)
  * Show pod resource limit/requests (#215, @bryanl)
  * Fixed error where capabilites for a plugin are not shown (#212, @bryanl)
  * Fixed bug where port forward states were active after pod is deleted (#209, @GuessWhoSamFoo)
  * 0.0.0.0 as listener/accepted will allow all hosts (#199, @wwitzel3)
  * Namespace drop down not refreshing on context change (#206, @bryanl)
  * Log errors when listing CRDs (#195, @wwitzel3)

### Enhancements
  * Use dynamic client when informer cache has not synced (#268, @bryanl)
  * Use d3-graphviz directly (remove graphdot-lib) (#248, @wwitzel3)
  * Convert actions to websockets and allow alerts from actions (#267, @bryanl)
  * Changed button label for removing port forward from "Remove" to "Stop port forward" (#226, @theneva)
  * Object store Get returns if object is found (#211, @bryanl)
  * Store for internal errors (#229, @wwitzel3)
  * Consolidated internal printers and tests (#42, @GuessWhoSamFoo)

### Documentation  
  * Added instructions to generate web assets (#252, @dotNomad)
  * Added hacking instructions to HACKING.md (#235, @wwitzel3)
  * Updated verbosity documentations (#149, @wwitzel3)
  * Added new resource examples (#220, @chelnak)