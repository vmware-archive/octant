import React from 'react'
import Section from './components/Section'
import './styles.scss'

export default function Navigation (props) {
  const { navigation } = props
  if (!navigation) return null
  const { sections = [] } = navigation
  return (
    <nav className='navigation--left'>
      {sections.map((section, i) => (
        <Section
          key={section.title}
          title={section.title}
          link={section.link}
          items={section.children}
        />
      ))}
    </nav>
  )
}
