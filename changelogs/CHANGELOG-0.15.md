## v0.15.0
#### 2020-08-11

### Download
 - https://github.com/vmware-tanzu/octant/releases/v0.15.0

### Highlights
  * Added printers for APIServices, MutatingWebhookConfigurations, and ValidatingWebhookConfigurations (#1045, @scothis)
  * Added stepper component (#1102, @wwitzel3, @nanaasiedu)
  * Added alert when unable to find valid local port (#1144, @ipsi)
  * Refactored JavaScript dashboard client to allow extending externally (#1218, @bryanl)
  * Added `client-max-recv-msg-size` flag to configure GRPC message size (#1231, @nodece)

### Bug Fixes
  * Fixed stepper component to send correct actions through websocket (#1102, @nanaasiedu)
  * Fixed inconsistent theming of editors (#1130, @scothis)
  * Fixed container dropdown losing focus on certain browsers (#1151, @GuessWhoSamFoo)
  * Fixed bug where kubeconfigs specified by multiple file paths were not watched (#1176, @GuessWhoSamFoo)
  * Fixed UI stutter when changing views (#1205, @mklanjsek)
  * Fixed bug to show correct CRD version when generating content (#1215, @bryanl)
  * Removed height from FlexLayoutSection (#1237, @GuessWhoSamFoo)
  * Fixed log components extending beyond frame in a plugin (#1245, @mklanjsek)

