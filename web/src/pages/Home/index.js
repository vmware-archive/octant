import React from 'react'
import Table from 'components/Table'
import './styles.scss'

function Home (props) {
  const { table } = props
  return (
    <div className='home'>
      <div className='main'>
        <Table data={table} />
      </div>
    </div>
  )
}

export default Home
