import { TextModel } from 'models/View'
import React from 'react'

interface Props {
  view: TextModel;
}

export default function({ view }: Props) {
  return <>{view.value}</>
}
