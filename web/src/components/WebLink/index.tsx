import { LinkModel } from 'models/Link'
import React from 'react'

interface Props {
    view: LinkModel;
}

export default function WebLink({view}: Props) {
    return <a href={view.ref}>{view.value}</a>
}
