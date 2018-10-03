import React from 'react'
import Select from 'react-select'
import './styles.scss'

export default function (props) {
  return (
    <Select
      className='dd-selector-container'
      classNamePrefix='dd-selector'
      {...props}
    />
  )
}
