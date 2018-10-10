import deployments from './_deployments'
import pods from './_pods'
import replicaSets from './_replicasets'
import workloads from './_workloads'
import ingresses from './_ingresses'
import services from './_services'
import discovery from './_discovery'

export default {
  'api/v1/content/overview/workloads/deployments': deployments,
  'api/v1/content/overview/workloads/pods': pods,
  'api/v1/content/overview/workloads/replica-sets': replicaSets,
  'api/v1/content/overview/workloads': workloads,
  'api/v1/content/overview/discovery-and-load-balancing/ingresses': ingresses,
  'api/v1/content/overview/discovery-and-load-balancing/services': services,
  'api/v1/content/overview/discovery-and-load-balancing': discovery
}
