import React from 'react'
import './styles.scss'

function Home (props) {
  const { summary, table } = props
  return (
    <div className='home'>
      <p>Welcome to the heptio ui-starter!</p>
      <div>
        {JSON.stringify(summary)}
      </div>
      <div>
        {JSON.stringify(table)}
      </div>
    </div>
  )
}

export default Home
