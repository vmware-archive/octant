import _ from 'lodash'
import React from 'react'

import NavigationHeader from './components/NavigationHeader'
import NavigationSubheader from './components/NavigationSubheader'
import './styles.scss'

interface Props {
  linkPath: NavigationSectionType[]
  childLinks: NavigationSectionType[]
  onNavChange: (NavigationPathLink) => void
}

export default function NavigationSection(props: Props) {
  const { linkPath: parentLinkPath, childLinks, onNavChange } = props
  return (
    <div className='navigation--left-section'>
      <NavigationHeader linkPath={parentLinkPath} onNavChange={onNavChange} />
      <ul className='navigation--left-items'>
        {_.map(childLinks, (link) => (
          <div key={link.title} className='navigation--left-item'>
            <NavigationSubheader
              linkPath={[...parentLinkPath, link]}
              childLinks={link.children}
              onNavChange={onNavChange}
            />
          </div>
        ))}
      </ul>
    </div>
  )
}
