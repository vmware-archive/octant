import React from 'react'
import QuickViewTitle from './components/QuickView/components/Title'
import { ResourceObject } from './schema'
import './quickview.scss'

interface Props {
  object: ResourceObject;
}

export default function({ object }: Props) {
  return (
    <div className='quickView'>
      <QuickViewTitle name={object.name} kind={object.kind} />
    </div>
  )
}
