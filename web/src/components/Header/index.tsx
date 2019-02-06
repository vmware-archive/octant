import './styles.scss'
import 'react-tagsinput/react-tagsinput.css'

import React, { Component } from 'react'
import Select from 'react-select'
import TagsInput from 'react-tagsinput'

interface Props {
  namespaceOptions: NamespaceOption[];
  namespace: string;
  namespaceValue: NamespaceOption;
  onNamespaceChange: (NamespaceOption) => void;
  resourceFilters: string[];
  onResourceFiltersChange: (filterTags: string[]) => void;
}

export default class extends Component<Props> {
  state = {
    tags: [],
  }

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
            <TagsInput
              inputProps={{ placeholder: 'Filter by label' }}
              value={resourceFilters}
              onChange={onResourceFiltersChange}
            />
          </div>
          <div className='header--context'>kubecontext</div>
        </div>
      </header>
    )
  }
}
