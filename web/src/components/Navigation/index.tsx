import React from 'react'
import _ from 'lodash'
import Selector from 'components/Selector'
import NavigationSection from './components/NavigationSection'
import './styles.scss'

interface Props {
  navSections: NavigationSectionType[];
  currentNavLinkPath: NavigationSectionType[];
  onNavChange: (NavigationPathLink) => void;
  namespaceOptions: NamespaceOption[];
  namespaceValue: NamespaceOption;
  onNamespaceChange: (NamespaceOption) => void;
}

export default function Navigation({
  navSections,
  currentNavLinkPath,
  onNavChange,

  namespaceOptions,
  namespaceValue,
  onNamespaceChange,
}: Props) {
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
