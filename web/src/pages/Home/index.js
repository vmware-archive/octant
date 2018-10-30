import React from 'react'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import Content from './components/Content'
import './styles.scss'

export default function Home ({ loading, contents, title }) {
  let mainContent
  if (loading) {
    mainContent = (
      <div className='loading-parent'>
        <Loading />
      </div>
    )
  } else if (contents && contents.length) {
    mainContent = contents.map((content, i) => (
      <div key={i} className='component--primary'>
        <Content content={content} />
      </div>
    ))
  } else {
    mainContent = (
      <div className='component--primary'>
        <h3 className='empty-content-text'>There is nothing to display here</h3>
      </div>
    )
  }

  return (
    <div className='home'>
      <Title title={title} />
      <div className='main'>{mainContent}</div>
    </div>
  )
}
