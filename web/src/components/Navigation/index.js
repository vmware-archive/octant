import React from 'react'
import Section from './components/Section'
import './styles.scss'

export default function Navigation (props) {
  const { navigation } = props
  if (!navigation) return null
  const { sections = [] } = navigation
  return (
    <nav className='navigation--left'>
      {sections.map(({ title, path, children }) => (
        <Section key={title} title={title} path={path} items={children} />
      ))}
    </nav>
  )
}
