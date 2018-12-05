import React from 'react'
import './styles.scss'

interface Props {
  name: string;
  kind: string;
}

export default function Title({ name, kind }: Props) {
  return (
    <div className='quickView-title'>
      <div className='name'>
        {name}
      </div>
      <div className='kind'>
        {kind}
      </div>
    </div>
  )
}
