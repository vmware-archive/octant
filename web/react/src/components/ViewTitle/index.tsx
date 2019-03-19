import { LinkModel, TextModel, TitleView } from 'models'
import React from 'react'
import { Link } from 'react-router-dom'

interface Props {
  parts: TitleView
}

export function ViewTitle(props: Props) {
  const { parts } = props

  if (!parts) {
    return null
  }

  const x = parts
    .map((part, i) => {
      switch (part.type) {
        case 'text':
          const text = part as TextModel
          return <span key={i}>{text.value}</span>
        default:
        case 'link':
          const link = part as LinkModel
          return <Link to={link.ref}>{link.value}</Link>
      }
    })
    .map((element, i) => {
      if (i < parts.length - 1) {
        return [
          element,
          <span className='component--title-separator' key={`sep${i}`}>
            &rsaquo;
          </span>,
        ]
      }
      return [element]
    })

  const titleParts = flatten(x)

  return <>{titleParts}</>
}

function flatten(a, depth = 1) {
  return a.reduce((flat, toFlatten) => {
    return flat.concat(Array.isArray(toFlatten) && depth - 1 ? toFlatten.flat(depth - 1) : toFlatten)
  }, [])
}
