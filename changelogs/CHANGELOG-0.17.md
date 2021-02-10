## v0.17.0
#### 2021-2-11

### Download
 - https://github.com/vmware-tanzu/octant/releases/v0.17.0

### Highlights
  * Added dropdowns to breadcrumbs (#1212, @mklanjsek)
  * Redesigned left navigation (#1353, @mklanjsek)
  * Added tooltips, support for different segment colors and thickness to donut charts (#1465, @mklanjsek)
  * Added Preferences to the bottom of the Vertical Navigation (#1498, @mklanjsek)
  * Added dropdown component (#1562, @mklanjsek)
  * Added pods table to node overview (#1773, @zparnold)
  * Added official language for approver status (#1874, @wwitzel3)
  * Added version skew policy (#1896, @GuessWhoSamFoo)

### Bug Fixes
  * Fixed issues with incorrect paths on Windows (#1696, @mklanjsek)
  * Fixed issue with broken links created from plugin (#1723, @mklanjsek)
  * Fixed streamer test flake (#1797, @GuessWhoSamFoo)

### Electron
  * Added preferences to bottom of the left nav (#1538, @mklanjsek)
  * Updated electron for osx arm64 (#1822, @akhenakh)
  * Added build pipeline using electron-builder (#1827, @GuessWhoSamFoo)
  * Changed to use random port when running as electron (#1852, @GuessWhoSamFoo)
  * Added menu to open log files (#1889, @wwitzel)
  * Added forward and back buttons to electron build (#1902, @GuessWhoSamFoo)
  * Consolidated preferences storage (#1957, @mklanjsek)
  * Added developer preferences (#1986, @mklanjsek)

### All Changes
  * Added docs section to reference.octant.dev (#1079, @wwitzel3)
  * Normalize markdown styles to match other text in datagrids (#1503, @scothis)
  * Updated dependabot to automatically vendor go module updates (#1524, @scothis)
  * Added multi key table sort and added reverse method (#1566, @GuessWhoSamFoo)
  * Refactored dash.Runner startup (#1676, @jamieklassen)
  * Show only preferred versions of CRDs instead of all versions (#1737, @bryanl)
  * Added support for windows shells and bash (#1749, @GuessWhoSamFoo)
