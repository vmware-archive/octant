import React from 'react'
import Section from './components/Section'
import './styles.scss'

export default function Summary (props) {
  const { data } = props
  const { sections, title = '' } = data
  return (
    <div className='summary--component'>
      <h2 className='summary-component-title'>{title}</h2>
      <hr />
      <div className='summary--component-sections'>
        {sections
          ? sections.map((section) => {
            const sectionTitle = section.title[0] === '_' ? '' : section.title
            return (
              <Section
                key={sectionTitle}
                title={sectionTitle}
                items={section.items}
              />
            )
          })
          : ''}
      </div>
    </div>
  )
}
