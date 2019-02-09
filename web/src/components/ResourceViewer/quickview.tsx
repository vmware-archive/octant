import React from 'react'

import QuickViewTitle from './components/QuickView/components/Title'
import './quickview.scss'
import { ResourceObject } from './schema'

interface Props {
  object: ResourceObject
}

export default function({ object }: Props) {
  return (
    <div className='quickView'>
      <QuickViewTitle name={object.name} kind={object.kind} />
    </div>
  )
}
