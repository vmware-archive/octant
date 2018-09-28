import React from 'react'
import Table from '../../components/Table'
import Summary from '../../components/Summary'

import './styles.scss'

function Home (props) {
  const { table, summary } = props
  return (
    <div className='home'>
      <div className='main'>
        <div className='component--primary'>
          <Table data={table} />
        </div>
        <div className='component--primary'>
          <Summary data={summary} />
        </div>
      </div>
    </div>
  )
}

export default Home
