import './styles.scss';

import React from 'react';

export default function ({ emptyContent }) {
  return (
    <div className='content-empty'>
      {emptyContent}
    </div>
  )
}
