import './styles.scss'

import _ from 'lodash'
import { GridModel } from 'models/View'
import React from 'react'
import GridLayout from 'react-grid-layout'
import { renderView } from 'views'

interface Props {
  view: GridModel;
}

export default function Grid({ view }: Props) {
  const panels = view.panels

  return (
    <GridLayout
      className='grid-layout'
      cols={24}
      rowHeight={25}
      width={1000}
      margin={[15, 10]}
      verticalCompact={true}
      compactType='vertical'
      autoSize={true}
    >
      {_.map(panels, (panel, i) => {
        const { position, content } = panel
        const dataGrid = { ...position, static: true }
        return (
          <div className='grid-layout-panel' key={i} data-grid={dataGrid}>
            {(() => {
              switch (content.type) {
                case 'labels':
                  return <div className='podtemplate-labels'>
                    <h3>{content.title}</h3>
                    {renderView(content)}
                  </div>
                default:
                return renderView(content)
              }

            })()}
          </div>
        )
      })}
    </GridLayout>
  )
}
