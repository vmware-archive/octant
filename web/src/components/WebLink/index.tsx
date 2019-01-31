import { LinkModel } from 'models/Link'
import React from 'react'
import { Link } from 'react-router-dom'

interface Props {
  view: LinkModel;
}

export default function WebLink({ view }: Props) {
    return <Link to={view.ref}>{view.value}</Link>
}
