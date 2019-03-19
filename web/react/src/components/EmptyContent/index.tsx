import React from 'react'

import './styles.scss'

interface Props {
  emptyContent: string
}

export default function({ emptyContent }: Props) {
  return <div className='content-empty'>{emptyContent}</div>
}
