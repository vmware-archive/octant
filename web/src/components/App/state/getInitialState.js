import Promise from 'promise'
import _ from 'lodash'
import { getNamespace, getNamespaces, getNavigation } from 'api'
import fetchContents from './fetchContents'

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
    return { loading: false, error: true }
  }

  const initialState = {}

  if (navigation) {
    initialState.navigation = navigation

    let currentNavLinkPath
    _.forEach(navigation.sections, (section) => {
      const linkPath = [section]
      if (section.path === currentPathname) {
        currentNavLinkPath = linkPath
        return false
      }
      _.forEach(section.children, (child) => {
        const childLinkPath = [...linkPath, child]
        if (child.path === currentPathname) {
          currentNavLinkPath = childLinkPath
          return false
        }
        _.forEach(child.children, (grandChild) => {
          const grandChildLinkPath = [...childLinkPath, grandChild]
          if (_.includes(currentPathname, grandChild.path)) {
            currentNavLinkPath = grandChildLinkPath
            return false
          }
        })
      })
    })

    if (currentNavLinkPath) {
      initialState.currentNavLinkPath = currentNavLinkPath
    }
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
      const contents = await fetchContents(currentPathname, option.value)
      _.assign(initialState, contents)
    }
  }

  return initialState
}
