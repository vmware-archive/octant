import Promise from 'promise'
import _ from 'lodash'
import { getNamespace, getNamespaces, getNavigation } from 'api'
import fetchContents from './fetchContents'

function getNavLinkPath (navigation, pathname) {
  let currentNavLinkPath

  if (navigation) {
    _.forEach(navigation.sections, (section) => {
      const linkPath = [section]
      if (section.path === pathname) {
        currentNavLinkPath = linkPath
        return false
      }
      _.forEach(section.children, (child) => {
        const childLinkPath = [...linkPath, child]
        if (child.path === pathname) {
          currentNavLinkPath = childLinkPath
          return false
        }
        _.forEach(child.children, (grandChild) => {
          const grandChildLinkPath = [...childLinkPath, grandChild]
          if (_.includes(pathname, grandChild.path)) {
            currentNavLinkPath = grandChildLinkPath
            return false
          }
        })
      })
    })
  }

  return currentNavLinkPath
}

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
      const contents = await fetchContents(currentPathname, option.value)
      _.assign(initialState, contents)
    }
  }

  return initialState
}
