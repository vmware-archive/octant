import './resource.scss'

import * as React from 'react'
import { createRef } from 'react'

import { ResourceObject } from './schema'

const resourceStyle = {
  height: '75px',
  width: '200px',
}

export interface Rect {
  top: number;
  left: number;
  height: number;
  width: number;
}

interface Props {
  name: string;
  object: ResourceObject;
  selected: boolean;
  setOffset(name: string, rect: Rect): void;
  updateSelected(name: string): void;
}

class Resource extends React.Component<Props, any> {
  private resourceRef = createRef<HTMLDivElement>()

  constructor(props: Props) {
    super(props)
    this.resourceRef = React.createRef()

    this.updateSelection = this.updateSelection.bind(this)
  }

  updateSelection(_: React.MouseEvent<HTMLDivElement>) {
    this.props.updateSelected(this.props.name)
  }

  componentDidMount() {
    if (this.resourceRef.current) {
      const current = this.resourceRef.current

      const rect: Rect = {
        top: current.offsetTop,
        left: current.offsetLeft,
        height: current.offsetHeight,
        width: current.offsetWidth,
      }

      this.props.setOffset(this.props.name, rect)
    }
  }

  render() {
    const object = this.props.object
    const selected = this.props.selected ? 'selected' : ''
    return (
      <div
        ref={this.resourceRef}
        className={`resource ${selected}`}
        style={resourceStyle}
        onClick={this.updateSelection}
      >
        <div className='resource-name'>{object.name}</div>
        <div className='resource-type'>
          {object.apiVersion} {object.kind}
        </div>
        <div className={`resource-status status-${object.status}`} />
      </div>
    )
  }
}

export default Resource
