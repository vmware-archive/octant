import deployments from './_deployments'
import pods from './_pods'
import replicaSets from './_replicasets'

export default {
  contents: [...deployments.contents, ...pods.contents, ...replicaSets.contents]
}
