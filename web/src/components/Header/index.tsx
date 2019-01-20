import './styles.scss'

import React from 'react'
import Select from 'react-select'

interface Props {
  namespaceOptions: NamespaceOption[];
  namespace: string;
  namespaceValue: NamespaceOption;
  onNamespaceChange: (NamespaceOption) => void;
}

class Header extends React.Component<Props> {
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
        <ul className='header--container'>
          <li className='header--logo'>
            <h1>Dev Dash</h1>
          </li>
          <li className='header--namespace'>
            <Select
              className='header--selector'
              classNamePrefix='header--selector'
              placeholder='Select namespace...'
              options={this.props.namespaceOptions}
              value={this.props.namespaceValue}
              onChange={this.props.onNamespaceChange}
            />
          </li>
          <li className='header--filter'>
            <input type='text' placeholder='label filter' />
          </li>
          <li className='header--context'>context display</li>
        </ul>
      </header>
    )
  }
}

export default Header
