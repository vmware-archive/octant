import { PortForwardModel } from 'models'
import React from 'react'

interface Props {
    view: PortForwardModel
}

export default function({view}: Props) {
    return <span>{view.text}</span>
}
