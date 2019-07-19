 - [v0.4.0](#v040)

## v0.4.0
#### 2019-07-19

### Download
- https://github.com/vmware/octant/releases/tag/v0.4.0

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
