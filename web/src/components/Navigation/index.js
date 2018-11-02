import React from 'react'
import _ from 'lodash'
import Selector from 'components/Selector'
import NavigationSection from './components/NavigationSection'
import './styles.scss'

export default function Navigation ({
  navSections,
  currentNavLinkPath,
  onNavChange,

  namespaceOptions,
  namespaceValue,
  onNamespaceChange
}) {
  return (
    <nav className='navigation--left'>
      <div className='navigation-namespace-selector'>
        <Selector
          placeholder='Select namespace...'
          options={namespaceOptions}
          value={namespaceValue}
          onChange={onNamespaceChange}
        />
      </div>
      {_.map(navSections, section => (
        <NavigationSection
          key={section.title}
          currentLinkPath={currentNavLinkPath}
          linkPath={[section]}
          childLinks={section.children}
          onNavChange={onNavChange}
        />
      ))}
    </nav>
  )
}
