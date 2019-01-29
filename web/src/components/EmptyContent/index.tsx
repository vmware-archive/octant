import './styles.scss'

import React from 'react'

interface Props {
  emptyContent: string;
}

export default function({ emptyContent }: Props) {
  return <div className='content-empty'>{emptyContent}</div>
}
