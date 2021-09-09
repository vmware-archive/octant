## v0.24.0

#### 2021-08-09

### Download

- https://github.com/vmware-tanzu/octant/releases/v0.24.0

### Highlights

- Added container image manifest for pods (#156, @mklanjsek)
- Added automatic reloading of Go plugins (#178, @wwitzel3)

### Bug Fixes

- Fixed status indicator for text and link components (#1688, @GuessWhoSamFoo)
- Fixed Go plugins spawning terminals under Windows (#2143, @wwitzel3)
- Fixed active tab not selected correctly when provided by plugins (#2809, @GuessWhoSamFoo)

### All Changes

- Added support for custom SVG icons (#1422, @ftovaro)
- Removed deprecated usage of `pb.Timestamp` (#2774, @GuressWhoSamFoo)
- Enabled Go plugins to apply yaml (#2775, @xtreme-vikram-yadav)
- Added better errors for JS plugin watcher (#2776, @xtreme-vikram-yadav)
- Added more granular logging for plugins (#2795, @lenriquez)
