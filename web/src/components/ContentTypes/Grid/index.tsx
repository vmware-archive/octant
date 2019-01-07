import React from 'react'
import _ from 'lodash'
import GridLayout from 'react-grid-layout'
import Summary, { ISummary } from 'components/ContentTypes/Summary'
import Table, { ITable } from 'components/ContentTypes/Table'
import Quadrant, { IQuadrant } from './components/Quadrant'
import Label, { ILabel } from './components/Label'
import './styles.scss'

type GridPanel = ContentType & {
  metadata: {
    type: 'panel';
  };
  config: {
    content: ContentType;
    position: GridPosition;
  }
}

export interface IGrid {
  metadata: {
    type: 'grid';
    title: string;
  };
  config: {
    panels: GridPanel[];
  };
}

interface Props {
  data: IGrid;
}

export default function Grid({ data }: Props) {
  const { config: { panels } } = data

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
      {
        _.map(panels, (panel, i) => {
          const { config: { position, content } } = panel
          const dataGrid = { ...position, static: true }
          return (
            <div className='grid-layout-panel' key={i} data-grid={dataGrid} >
              {
                (() => {
                  switch (content.metadata.type) {
                    case 'quadrant':
                      return <Quadrant data={content as IQuadrant} />
                    case 'label':
                      return <Label data={content as ILabel} />
                    case 'summary':
                      return <Summary data={content as ISummary} />
                    case 'table':
                      return <Table data={content as ITable} />
                    default:
                      return `unknown content type [${content.metadata.type}]`
                  }
                })()
              }
            </div>
          )
        })
      }
    </GridLayout>
  )
}
