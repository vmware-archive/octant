import React, { Component } from 'react'
import Select from 'react-select'

import Filter from './components/Filter'
import './styles.scss'

interface Props {
  namespaceOptions: NamespaceOption[]
  namespace: string
  namespaceValue: NamespaceOption
  onNamespaceChange: (NamespaceOption) => void
  resourceFilters: string[]
  onResourceFiltersChange: (filterTags: string[]) => void
}

export default class extends Component<Props> {
  state = {
    tags: [],
  }

  render() {
    const { namespaceOptions, namespaceValue, onNamespaceChange, resourceFilters, onResourceFiltersChange } = this.props
    return (
      <header>
        <div className='header--container'>
          <div className='header--logo'>
            <h1>Sugarloaf</h1>
          </div>
          <div className='header--namespace'>
            <Select
              className='header--selector'
              classNamePrefix='header--selector'
              placeholder='Select namespace...'
              options={namespaceOptions}
              value={namespaceValue}
              onChange={onNamespaceChange}
            />
          </div>
          <div className='header--filter'>
            <Filter
              resourceFilters={resourceFilters}
              onResourceFiltersChange={onResourceFiltersChange}
            />
          </div>
        </div>
      </header>
    )
  }
}
