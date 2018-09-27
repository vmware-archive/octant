import React from 'react'
import Section from './components/Section'
import NavigationData from './_mock_data'

import './styles.scss'

const navData = NavigationData.sections

function Navigation () {
  return (
    <nav className='navigation--left'>
      {navData.map(section => (
        <Section
          name={section.name}
          link={section.link}
          items={section.children}
          key={section.link}
        />
      ))}
    </nav>
  )
}

export default Navigation
