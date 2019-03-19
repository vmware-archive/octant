import _ from 'lodash'
import { AnnotationsModel } from 'models'
import React from 'react'

import './styles.scss'

interface Props {
  view: AnnotationsModel
}

export default function({ view: { annotations } }: Props) {
  return (
    <div className='content-annotations'>
      {_.map(annotations, (value, key) => {
        return (
          <div key={key} className='content-annotation'>
            <span className='content-annotation-key'>{key}</span>: {value}
          </div>
        )
      })}
    </div>
  )
}
