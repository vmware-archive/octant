import Promise from 'promise'
import _ from 'lodash'
import { getNamespace, getNamespaces, getNavigation } from 'api'
import getNavLinkPath from './getNavLinkPath'

export default async function (currentPathname) {
  let navigation,
    namespaces,
    namespace
  try {
    [navigation, namespaces, namespace] = await Promise.all([
      getNavigation(),
      getNamespaces(),
      getNamespace()
    ])
  } catch (e) {
    return { isLoading: false, hasError: true }
  }

  const initialState = {
    navigation,
    currentNavLinkPath: getNavLinkPath(navigation, currentPathname)
  }

  if (namespaces && namespaces.namespaces && namespaces.namespaces.length) {
    initialState.namespaceOptions = namespaces.namespaces.map(ns => ({
      label: ns,
      value: ns
    }))
  }

  if (namespace && initialState.namespaceOptions.length) {
    const option = _.find(initialState.namespaceOptions, {
      value: namespace.namespace
    })
    if (option) {
      initialState.namespaceOption = option
    }
  }

  return initialState
}
