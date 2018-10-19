import React from 'react'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import ContentSwitcher from './components/ContentSwitcher'
import './styles.scss'

export default function Home ({ loading, contents, title }) {
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
      <Title title={title} />
      <div className='main'>
        {mainContent}
      </div>
    </div>
  )
}
