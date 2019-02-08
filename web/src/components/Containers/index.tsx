import { ContainersModel } from 'models/View'
import React from 'react'
import './styles.scss'

interface Props {
  view: ContainersModel;
}

export default function({ view }: Props) {
  return (
    <>
      {view.containerDefs.map((containerDef, index) => {
        return (
          <>
            <div key={index} className='table-containers'>
              <div className='table-containers-name'>
                {containerDef.name}
              </div>
              <div className='table-containers-image'>
                {containerDef.image}
              </div>
            </div>
          </>
        )
      })}
    </>
  )
}
