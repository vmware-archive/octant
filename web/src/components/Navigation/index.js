import React from 'react'
import Selector from 'components/Selector'
import Section from './components/Section'
import './styles.scss'

export default function Navigation ({
  navigation,
  namespaceOptions,
  namespaceValue,
  onNamespaceChange
}) {
  if (!navigation) return null
  const { sections = [] } = navigation
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
      {sections.map(({ title, path, children }) => (
        <Section key={title} title={title} path={path} items={children} />
      ))}
    </nav>
  )
}
