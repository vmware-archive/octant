import React from 'react'
import { Link } from 'react-router-dom'
import './styles.scss'

interface Props {
  params: LinkContentType;
}

export default function Item(props: Props) {
  const {
    metadata: { title },
    config: { value, ref },
  } = props.params
  return (
    <div className='summary--data summary--data-link'>
      <div className='summary--data-key'>{title}</div>
      <div className='summary--data-link'>
        <Link className='link--gray' to={ref}>
          {value || ref}
        </Link>
      </div>
    </div>
  )
}
