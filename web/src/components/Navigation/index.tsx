import './styles.scss'

import _ from 'lodash'
import React from 'react'

import NavigationSection from './components/NavigationSection'

interface Props {
  navSections: NavigationSectionType[];
  currentNavLinkPath: NavigationSectionType[];
  onNavChange: (NavigationPathLink) => void;
}

export default function Navigation({
  navSections,
  currentNavLinkPath,
  onNavChange,
}: Props) {
  return (
    <nav className='navigation--left'>

      {_.map(navSections, (section) => (
        <NavigationSection
          key={section.title}
          linkPath={[section]}
          childLinks={section.children}
          onNavChange={onNavChange}
        />
      ))}
    </nav>
  )
}
