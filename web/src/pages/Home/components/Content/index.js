import React from 'react'
import Table from 'components/ContentTypes/Table'
import Summary from 'components/ContentTypes/Summary'

export default function ({ content }) {
  const { type } = content
  if (type === 'table') {
    return <Table data={content} />
  }
  if (type === 'summary') {
    return <Summary data={content} />
  }
  return <div>Can not render content type</div>
}
