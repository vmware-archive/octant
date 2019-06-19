  - [v0.2.1](#v021)
  - [v0.2.0](#v020)

## v0.2.1
#### 2019-04-18

### Download
- https://github.com/vmware/octant/releases/tag/v0.2.1

### Bug Fixes / Other Changes:
  * Generated more detailed error output when a describer tab function fails. Also, if a tab function fails, do not stop processing. (#737, @bryanl)
  * A user can clear all filters from the UI now (#753, @mdaverde)
  * You can also disable timestamps in your container logs (#753, @mdaverde)
  * Fixed regression that doesn't allow event lists with node events to render (#740, @bryanl)
  * Used a mutex when sorting table rows because it is possible the table will change while it is being sorted. (#740, bryanl)
  * Don't try to link to nodes because dash doesn't support them (#740, bryanl)
  * Don't try to sort event tables by column name (it doesn't exist) (#733, @bryanl)
  * Dashboard service should only listen on localhost (#735, @bryanl)
  * Checked the host header. If it isn't localhost, or 127.0.0.1, return with forbidden status #738 (#738, @bryanl)
  * Prevented progress bar from pushing down the rest of the page (#746, @mdaverde)

## v0.2.0
#### 2019-04-15

### Download
- https://github.com/vmware/octant/releases/tag/v0.2.0

### Highlights
- Object viewer with the ability to switch between namespaces.
- Resource viewer, a relationship graph.
- Container logs
- Port forwarding setup via UI
- Live updating streamed from the cluster as it changes.
