import React from 'react'
import Section from './components/Section'
import './styles.scss'

export default function Navigation (props) {
  const { navigation } = props
  return (
    <nav className='navigation--left'>
      {navigation.map(section => (
        <Section
          title={section.title}
          link={section.link}
          items={section.children}
          key={section.link}
        />
      ))}
    </nav>
  )
}
