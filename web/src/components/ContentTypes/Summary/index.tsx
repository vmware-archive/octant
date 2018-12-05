import React from 'react'
import _ from 'lodash'
import Section from './components/Section'
import './styles.scss'

interface Props {
  data: ContentSummary,
}

export default function Summary(props: Props) {
  const { data } = props
  const { sections, title = '' } = data
  return (
    <div className='summary--component'>
      <h2 className='summary-component-title'>{title}</h2>
      <hr />
      <div className='summary--component-sections'>
        {_.map(sections, (section) => {
          const sectionTitle = section.title[0] === '_' ? '' : section.title
          return (
            <Section
              key={sectionTitle}
              title={sectionTitle}
              items={section.items}
            />
          )
        })}
      </div>
    </div>
  )
}
