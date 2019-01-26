import { ContainersModel } from 'models/View'
import React from 'react'

interface Props {
  view: ContainersModel;
}

export default function({ view }: Props) {
  return (
    <>
      {view.containerDefs.map((containerDef, index) => {
        return (
          <div key={index} className='containerdef'>
            {containerDef.name}
            =>
            {containerDef.image}
          </div>
        )
      })}
    </>
  )
}
