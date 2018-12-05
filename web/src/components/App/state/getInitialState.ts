import PromisePolyfill from 'promise'
import _ from 'lodash'
import { getNamespace, getNamespaces, getNavigation } from 'api'
import getNavLinkPath from './getNavLinkPath'

interface InitialState {
  isLoading?: boolean;
  hasError?: boolean;
  navigation?: { sections: NavigationSectionType[] };
  currentNavLinkPath?: NavigationSectionType[];
  namespaceOption?: NamespaceOption;
  namespaceOptions?: NamespaceOption[];
}

export default async function(currentPathname): Promise<InitialState> {
  let navigation, namespaces, namespace
  try {
    [navigation, namespaces, namespace] = await PromisePolyfill.all([
      getNavigation(),
      getNamespaces(),
      getNamespace(),
    ])
  } catch (e) {
    return { isLoading: false, hasError: true }
  }

  const initialState: InitialState = {
    navigation,
    currentNavLinkPath: getNavLinkPath(navigation, currentPathname),
  }

  if (namespaces && namespaces.namespaces && namespaces.namespaces.length) {
    initialState.namespaceOptions = namespaces.namespaces.map((ns) => ({
      label: ns,
      value: ns,
    }))
  }

  if (namespace && initialState.namespaceOptions.length) {
    const option: NamespaceOption = _.find(initialState.namespaceOptions, {
      value: namespace.namespace as string,
    })
    if (option) {
      initialState.namespaceOption = option
    }
  }

  return initialState
}
