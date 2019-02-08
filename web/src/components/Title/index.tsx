import React from 'react'

import './styles.scss'

export default function({ title }) {
  if (!title) return null
  return (
    <div className='component--title'>
      <h2>{title}</h2>
    </div>
  )
}
