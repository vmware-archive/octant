- [v0.10.1](#v0101)
- [v0.10.0](#v0100)

## v0.10.1
#### 2020-01-30

### Download
 - https://github.com/vmware/octant/releases/tag/v0.10.1

### All Changes
  * Fixed error type assertion for AccessError. (#562, @wwitzel3)
  * Fixed Workload navigation on Windows (#568, @wwitzel3)
  * Added status column to pod conditions (#569, @bryanl)
  * Fixed workload view failure if pod metrics are not available for a pod (#572, @bryanl)
  * Fixed typo on GO111MODULE in build script's goInstall() (#574, @ilayaperumalg)
  * Added a backoff strategy for objectstore access (#579, @wwitzel3)
  * Removed extra pod log view (#585, @bryanl)
  * Fixed CRD watcher race condition (#588, @GuessWhoSamFoo)
  * Added cluster client reload when kubeconfig changes (#591, @bryanl)
  * Disabled CircleCI builds (#603, @bryanl)

## v0.10.0
#### 2020-01-22

### Download
 - https://github.com/vmware/octant/releases/tag/v0.10.0

### Highlights
 - Check-in the generated mocks used for unit testing (#427, @antoninbas)
 - Added keyboard shortcuts for navigation (#434, @alexbrand)
 - Updated octant favicon to logo (#462, @bryanl)
 - Added workload module with support for metrics API to show resource usage for pods (#480, @bryanl)
 - Added terminal component and ability to execute a command against a container (#488, @wwitzel3)
 - Added support for v1 and v1beta1 CRD API (#490, @bryanl)
 - Updated new community meeting time on website (#535, @jonasrosland)
 - Moved object metadata to its own tab (#542, @bryanl)

### All Changes
  * Moved internal/dash to pkg/dash (#396, @wwitzel3)
  * Sorted events in a more stable manner (#404, @bryanl)
  * Fixed bug where container ports used for active port fowards used an incorrect state id (#409, @GuessWhoSamFoo)
  * Added iframe component (#431, @GarySmith)
  * Updated service printer to show service port instead of target port (#435, @alexbrand)
  * Created slider component (#440, @bryanl)
  * Fixed site accessibility and updated gem dependencies (#444, @SDBrett)
  * Updated documentation and add clarity around running and using Octant (#450, @wwitzel3)
  * Added uninstall instructions to plugin documentation (#465, @danielhelfand)
  * Removed typescript warnings in content switcher (#476, @bryanl)
  * Removed klog messages unless asked for explicitly (#529, @bryanl)
  * Moved port forward state management to backend (#543, @GuessWhoSamFoo)
  * Sorted container env vars and show summary with missing configmap references (#545, @GuessWhoSamFoo)
  * Updated resource viewer to show workload when it is selected (#551, @bryanl)
  * Updated documentation to use build.go instead of make (#559, @GarySmith)
