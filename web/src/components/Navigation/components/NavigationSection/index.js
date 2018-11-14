import React from 'react'
import _ from 'lodash'
import NavigationSubheader from './components/NavigationSubheader'
import NavigationHeader from './components/NavigationHeader'
import './styles.scss'

export default function NavigationSection (props) {
  const {
    currentLinkPath,
    linkPath: parentLinkPath,
    childLinks,
    onNavChange
  } = props
  return (
    <div className='navigation--left-section'>
      <NavigationHeader
        currentLinkPath={currentLinkPath}
        linkPath={parentLinkPath}
        onNavChange={onNavChange}
      />
      <ul className='navigation--left-items'>
        {_.map(childLinks, link => (
          <div key={link.title} className='navigation--left-item'>
            <NavigationSubheader
              currentLinkPath={currentLinkPath}
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
