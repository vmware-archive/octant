import React from 'react'
import _ from 'lodash'
import Selector from 'components/Selector'
import Section from './components/Section'
import './styles.scss'

export default function Navigation ({
  navigation,
  namespaceOptions,
  namespaceValue,
  onNamespaceChange
}) {
  return (
    <nav className='navigation--left'>
      {navigation && (
        <React.Fragment>
          <div className='navigation-namespace-selector'>
            <Selector
              placeholder='Select namespace...'
              options={namespaceOptions}
              value={namespaceValue}
              onChange={onNamespaceChange}
            />
          </div>
          {_.map(navigation.sections, ({ title, path, children }) => (
            <Section key={title} title={title} path={path} items={children} />
          ))}
        </React.Fragment>
      )}
    </nav>
  )
}
