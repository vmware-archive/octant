import React from 'react'

import { LogEntry } from './log'
import './styles.scss'

interface Props {
  log: LogEntry
}

export default function Message({ log }: Props) {
  return (
    <>
      <span className='logs--message-timestamp'>{log.timestamp}</span>
      <span className='logs--message-text'>{log.message}</span>
    </>
  )
}
