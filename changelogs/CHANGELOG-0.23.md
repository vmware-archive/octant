## v0.23.0

#### 2021-08-02

### Download

- https://github.com/vmware-tanzu/octant/releases/v0.23.0

### Highlights

- Go Plugin API now can use `SendEvent`.
- Custom Resources will show conditions table if conditions are found.
- Resources now show their status as well as phase.
- Lots of improvements to object status.

### Bug Fixes

- Fixed issue with --namespace-list not adding namespaces list. (#2259, @wwitzel3)
- Fixed panic when viewing objects with no owner reference in applications view (#2650, @GuessWhoSamFoo)

### All Changes

- Added termination threshold (#1408, @lenriquez)
- Added `PVC` to object status (#2641, @ftovaro)
- Added Octant route history dropdown and updated page title (#2580, @xtreme-vikram-yadav)
- Changed node view to use object tables (#2665, @GuessWhoSamFoo)
- Added conditions table automatically to overview summary. (#489, @wwitzel3)
- Added `SendEvent` for go plugins (#2691, @GuessWhoSamFoo)
- Added status column to Pod list (#201, @ftovaro)
- Fixed issue with --namespace-list not adding namespaces list. (#2259, @wwitzel3)
- Fixed panic when viewing objects with no owner reference in applications view (#2650, @GuessWhoSamFoo)
- Fixed `ObjectStatusResponse` to update resource viewer and object tables (#2651, @GuessWhoSamFoo)
- Fixed issues preventing complex custom resource viewer graphs from being displayed in TypeScript-based plugins (#2669 @liamrathke)
- Fixed form select not sending changes on Submit (#2623, @lenriquez)
