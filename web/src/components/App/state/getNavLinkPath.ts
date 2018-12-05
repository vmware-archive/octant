import _ from 'lodash'

export default function getNavLinkPath(navigation: Navigation, pathname: string): NavigationSectionType[] {
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
