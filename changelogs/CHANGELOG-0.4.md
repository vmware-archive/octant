 - [v0.4.1](#v041)
 - [v0.4.0](#v040)

## v0.4.1
#### 2019-07-29

### Download
 - https://github.com/vmware-tanzu/octant/releases/tag/v0.4.1

### Changes
  * Upgrade to client-go 1.12 and k8s api 1.15 (#77, @wwitzel3)
  * Added namespace and portforward buttons to port forward list (#84, @guesswhosamfoo)
  * Fix showing multiple containers within a pod in summary view (#88, @GuessWhoSamFoo)
  * Preserve resource view during a namespace switch (#80, @mdaverde)
  * fix an issue where a bad CRD would prevent the UI from loading (#76, @wwitzel3)
  * Add octant flag for specifying klog verbosity level (#52, @bryanl)
  * Tune rate limiter for client-go's rest client (#51, @bryanl)
  * Rename sample plugin to octant-sample-plugin (#50, @bryanl)


## v0.4.0
#### 2019-07-19

### Download
- https://github.com/vmware-tanzu/octant/releases/tag/v0.4.0

### Highlights
- Plugin helper service to provide a method for plugin authors to initialize and implement plugins.
- Resource viewer now uses [Cytoscape.js](http://js.cytoscape.org/)
- Namespaces can be searched in addition the dropdown.
- Plugins can act as modules to server content and provide navigation

### All Changes:

  * Created plugin helper service (#17, @bryanl)
  * Sorted pods by name when printing in a list (#18, @bryanl)
  * Converted resource viewer to cytoscape JS (#7, @bryanl)
  * Added support for a card component with actions (#13, @bryanl)
  * Increased performacne by creating informer factories on demand (#10, @wwitzel3)
  * Add module support to plugins (#9, @bryanl)
  * Allow user to search against list of namespaces (#1, @mdaverde)
