import './graph.scss'

import * as React from 'react'

import Connections from './connections'
import Resource, { Rect } from './resource'
import { Schema } from './schema'

interface Props {
  rows: Array<Set<string>>;
  schema: Schema;
  setCurrentResource(name: string): void;
}

export interface State {
  offsets: {
    [key: string]: Rect;
  };
  schema: Schema;
  selected: string;
}

class Graph extends React.Component<Props, State> {
  static getDerivedStateFromProps(props: Props, state: State) {
    if (JSON.stringify(props.schema) !== JSON.stringify(state.schema)) {
      return {
        offsets: {},
        schema: props.schema,
      }
    }
    return null
  }

  constructor(props: Props) {
    super(props)
    this.setOffset = this.setOffset.bind(this)
    this.state = {
      offsets: {},
      schema: props.schema,
      selected: props.schema.selected,
    }

    this.updateSelected = this.updateSelected.bind(this)
  }

  setOffset(name: string, rect: Rect) {
    const offsets = this.state.offsets
    offsets[name] = rect
    this.setState({ offsets })
  }

  updateSelected(name: string) {
    this.props.setCurrentResource(name)
    this.setState({ selected: name })
  }

  render() {
    const rows = this.props.rows.map((row, rowID) => {
      const columns = [...row].map((name, columnID) => {
        const object = this.props.schema.objects[name]
        return (
          <Resource
            key={columnID}
            name={name}
            object={object}
            selected={this.state.selected === name}
            updateSelected={this.updateSelected}
            setOffset={this.setOffset}
          />
        )
      })
      return (
        <div key={rowID} className='row'>
          {columns}
        </div>
      )
    })

    return (
      <div className='graph'>
        <Connections offsets={this.state.offsets} adjacencyList={this.state.schema.adjacencyList} />
        {rows}
      </div>
    )
  }
}

export default Graph
