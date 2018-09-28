import React from 'react'
import Section from './components/Section'
import SummaryData from './_mock_data'

import './styles.scss'

export default function Summary () {
  const { title, sections } = SummaryData
  return (
    <div className='summary--component'>
      <h2 className='summary-component-title'>{title}</h2>
      <hr />
      {sections.map((section) => {
        const sectionTitle = section.type[0] === '_' ? '' : section.type
        return <Section title={sectionTitle} data={section.data} />
      })}
    </div>
  )
}
