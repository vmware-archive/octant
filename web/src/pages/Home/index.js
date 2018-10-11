import React from 'react'
import Loading from 'components/Icons/Loading'
import ContentSwitcher from './components/ContentSwitcher'
import './styles.scss'

export default function Home ({ loading, contents }) {
  let mainContent = <div>No resources</div>
  if (loading) {
    mainContent = (
      <div className='loading-parent'>
        <Loading />
      </div>
    )
  } else if (contents && contents.length) {
    mainContent = contents.map((content, i) => (
      <div key={i} className='component--primary'>
        <ContentSwitcher content={content} />
      </div>
    ))
  }

  return (
    <div className='home'>
      <div className='main'>{mainContent}</div>
    </div>
  )
}
