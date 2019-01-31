import React, { Component } from 'react'
import Select from 'react-select'
import './styles.scss'

interface Props {
  namespaceOptions: NamespaceOption[];
  namespace: string;
  namespaceValue: NamespaceOption;
  onNamespaceChange: (NamespaceOption) => void;
}

export default class extends Component<Props> {
  namespaces() {
    return this.props.namespaceOptions.map((option, i) => {
      return (
        <option key={i} value={option.value}>
          {option.label}
        </option>
      )
    })
  }

  render() {
    return (
      <header>
        <div className='header--container'>
          <div className='header--logo'>
            <h1>Dev Dash</h1>
          </div>
          <div className='header--namespace'>
            <Select
              className='header--selector'
              classNamePrefix='header--selector'
              placeholder='Select namespace...'
              options={this.props.namespaceOptions}
              value={this.props.namespaceValue}
              onChange={this.props.onNamespaceChange}
            />
          </div>
          <div className='header--filter'>
            <input className='header-filter-input' type='text' placeholder='Filter by label' />
          </div>
          <div className='header--context'>kubecontext</div>
        </div>
      </header>
    )
  }
}
