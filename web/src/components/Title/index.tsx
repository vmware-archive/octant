import { ViewTitle } from 'components/ViewTitle'
import { TitleView } from 'models'
import React from 'react'

import './styles.scss'

interface Props {
  parts: TitleView
}

export default function({ parts }: Props) {
  return (
    <div className='component--title'>
      <h2>
        <ViewTitle parts={parts} />
      </h2>
    </div>
  )
}
