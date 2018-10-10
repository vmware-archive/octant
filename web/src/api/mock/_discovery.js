import ingresses from './_ingresses'
import services from './_services'

export default {
  contents: [...ingresses.contents, ...services.contents]
}
