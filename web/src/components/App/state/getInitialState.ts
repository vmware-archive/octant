import { getNamespace, getNamespaces, getNavigation } from 'api'
import _ from 'lodash'
import PromisePolyfill from 'promise'

import getNavLinkPath from './getNavLinkPath'

interface InitialState {
  isLoading?: boolean
  hasError?: boolean
  navigation?: { sections: NavigationSectionType[] }
  currentNavLinkPath?: NavigationSectionType[]
  namespaceOption?: NamespaceOption
  namespaceOptions?: NamespaceOption[]
}

const namespaceRe = /\/content\/overview\/namespace\/([A-Za-z-]+)\//ig

export default async function(currentPathname): Promise<InitialState> {
  let navigation, namespaces, namespace: string

  try {
    // tslint:disable-next-line
    ;[navigation, namespaces] = await PromisePolyfill.all([getNavigation(), getNamespaces()])
  } catch (e) {
    return { isLoading: false, hasError: true }
  }

  const initialState: InitialState = {
    navigation,
    currentNavLinkPath: getNavLinkPath(navigation, currentPathname),
  }

  const matches = namespaceRe.exec(currentPathname)
  if (matches && matches.length > 1) {
    namespace = matches[1]
  }

  if (namespaces && namespaces.namespaces && namespaces.namespaces.length) {
    initialState.namespaceOptions = namespaces.namespaces.map((ns) => ({
      label: ns,
      value: ns,
    }))
  }

  const option: NamespaceOption = _.find(initialState.namespaceOptions, {
    value: namespace,
  })

  initialState.namespaceOption = option ? option : { label: 'default', value: 'default' }

  return initialState
}
