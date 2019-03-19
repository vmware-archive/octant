import React from 'react'
import ReactTooltip from 'react-tooltip'

import './styles.scss'

export function Tooltip(props) {
  return (
    <ReactTooltip id={props.id} offset={{ right: 3 }} effect='solid' place='right' className='viewerTooltip'>
      {props.children}
    </ReactTooltip>
  )
}
