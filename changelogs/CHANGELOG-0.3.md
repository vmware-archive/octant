 - [v0.3.0](#v030)

## v0.3.0
#### 2019-06-17

### Download
- https://github.com/vmware/octant/releases/tag/v0.3.0 

### Highlights
- Speed and UX improvements.
- Better error handling.
- Ability to switch contexts from the UI.
- Documentation for Plugins

### All Changes:
  * Fix resource viewer to make it stop sending invalid object graphs (#827, @bryanl)
  * Improve UX by adding loading and error pages. (#829, @wwitzel3)
  * Allow sorting sorting of link compnents (#818, @bryanl)
  * Add support for showing current context to frontend (#832, @bryanl)
  * Support cluster level custom resources (#817, @bryanl)
  * Fix panic when trying to get container statuses from unscheduled pods (#797, @bryanl)
  * Show service account from workload view (#692, @bryanl)
  * Added `OCTANT_PLUGIN_PATH` environment variable to take list of paths (#790, @GuessWhoSamFoo)
  * Added api endpoint for registered plugins (#788, @GuessWhoSamFoo)
  * Added clean command to Makefile for generated mock files. (#765, @GuessWhoSamFoo)
  * Added links to Service Accounts from Role Binding (#687, @GuessWhoSamFoo)
  * Resource Viewer now shows a loading message for graph generation that takes longer than 750ms (#585, @wwitzel3)
  * Check access per resource to provide a better user experience for non-admins. (#774, @wwitzel3)
