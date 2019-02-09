import React from 'react'
import { Link as RouterLink } from 'react-router-dom'

import './styles.scss'

interface Props {
  params: LinkContentType
}

export default function Link({ params }: Props) {
  const {
    metadata: { title },
    config: { value, ref },
  } = params
  return (
    <div className='summary--data summary--data-link'>
      <div className='summary--data-key' data-test='title'>
        {title}
      </div>
      <div className='summary--data-link'>
        <RouterLink className='link--gray' to={ref}>
          {value || ref}
        </RouterLink>
      </div>
    </div>
  )
}
