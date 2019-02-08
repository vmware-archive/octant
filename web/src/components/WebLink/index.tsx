import { LinkModel } from 'models'
import React from 'react'
import { Link } from 'react-router-dom'

import './styles.scss'

interface Props {
  view: LinkModel
}

export default function WebLink({ view }: Props) {
  return (
    <Link className='web-link' to={view.ref}>
      {view.value}
    </Link>
  )
}
