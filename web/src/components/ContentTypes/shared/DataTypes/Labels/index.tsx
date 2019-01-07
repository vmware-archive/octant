import React from 'react'
import Labels from 'components/ContentTypes/shared/Labels'

interface Props {
  params: LabelsContentType;
}

export default function Item(props: Props) {
  const { params } = props
  const {
    metadata: { title },
    config: { labels },
  } = params
  return (
    <div className='summary--data'>
      {
        title && (
          <div className='summary--data-key'>{title}</div>
        )
      }
      <div className='summary--data-labels'>
        <Labels labels={labels} />
      </div>
    </div>
  )
}
