import { ViewTitle } from 'components/ViewTitle'
import { FlexLayoutModel } from 'models'
import Content from 'pages/Overview/components/Content'
import React from 'react'
import { renderView } from 'views'

import './styles.scss'

interface Props {
  view: FlexLayoutModel
}

export default function FlexLayout({ view }: Props) {
  const sections = view.sections.map((section, sectionIndex) => {
    const items = section.map((item, itemIndex) => {
      const style = {
        // flexBasis: `${40 * item.width}px`,
        flexBasis: `${item.width / 24 * 100 - 5}%`,
      }

      return (
        <div key={itemIndex} className='flexLayout--section-item' style={style}>
          {(() => {
            switch (item.view.type) {
              case 'labels':
                return (
                  <div className='podtemplate-labels'>
                    <h3>
                      <ViewTitle parts={item.view.title} />
                    </h3>
                    {renderView(item.view)}
                  </div>
                )
              default:
                return renderView(item.view)
            }
          })()}
        </div>
      )
    })

    return (
      <div key={sectionIndex} className='flexLayout--section'>
        {items}
      </div>
    )
  })

  return <div className='flexLayout'>{sections}</div>
}
