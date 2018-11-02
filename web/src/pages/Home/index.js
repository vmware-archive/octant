import React from 'react'
import cx from 'classnames'
import Loading from 'components/Icons/Loading'
import Title from 'components/Title'
import Content from './components/Content'
import './styles.scss'

export default function Home ({
  loading, contents, title, error
}) {
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
    const cnames = cx({
      'content-text': true,
      'empty-content-text': error === false,
      'error-content-text': error === true
    })
    mainContent = (
      <div className='component--primary'>
        <h3 className={cnames}>
          {error === true
            ? "Oops, something's not right, try again."
            : 'There is nothing to display here'}
        </h3>
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
