import React from 'react'
import './styles.scss'

export default function ({ title }) {
  return (
    <div className='content-empty'>
      This namespace does not have any {title}.
    </div>
  )
}
