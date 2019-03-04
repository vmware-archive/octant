import React, { Component } from 'react'
import _ from 'lodash'
import DownArrow from './components/DownArrow'

import './styles.scss'

interface Props {
  resourceFilters: string[]
  onResourceFiltersChange: (filterTags: string[]) => void
}

interface State {
  inputValue: string
  showFilters: boolean
}

export default class Filter extends Component<Props, State> {
  state: State = {
    inputValue: '',
    showFilters: false,
  }

  private node = React.createRef<HTMLDivElement>()

  closeFilters = (evt) => {
    if (!document.body.contains(this.node.current)) return
    if (!document.body.contains(evt.target)) return

    if (this.node.current.contains(evt.target)) {
      return
    }

    this.setState({ showFilters: false })
  };

  componentDidMount() {
    document.addEventListener('click', this.closeFilters, false)
  }

  componentWillUnmount(): void {
    document.removeEventListener('click', this.closeFilters, false)
  }

  render() {
    const { resourceFilters, onResourceFiltersChange } = this.props;
    const { inputValue, showFilters } = this.state
    return (
      <div className='header-filter-component'  ref={this.node}>
        <div className='header-filter-input-wrapper'>
          <div className='header-filter-input'>
            <input
              type='text'
              placeholder='Filter by label'
              value={inputValue}
              onChange={(evt) => {
                this.setState({ inputValue: evt.currentTarget.value })
              }}
              onKeyPress={(evt) => {
                const { value } = evt.currentTarget
                if (evt.key === 'Enter' && value) {
                  const newResourceFilters = [...resourceFilters, value]
                  onResourceFiltersChange(newResourceFilters)
                  this.setState({ inputValue: '', showFilters: true })
                }
              }}
            />
          </div>
          <div
            className='header-filter-downarrow'
            onClick={() => {
              const { showFilters } = this.state
              this.setState({ showFilters: !showFilters })
            }}
          >
            <DownArrow />
          </div>
        </div>
        {
          showFilters ? (
            <div className='header-filter-tagslist'>
              {
                resourceFilters && resourceFilters.length ? (
                    _.map(resourceFilters, (filter, index) => {
                      return (
                        <div key={index} className='header-filter-tagslist-row'>
                          <div className='header-filter-tagslist-tag'>
                            {filter}
                            <a
                              className='header-filter-tagslist-remove'
                              onClick={() => {
                                const newResourceFilters = [...resourceFilters]
                                _.pullAt(newResourceFilters, [index])
                                onResourceFiltersChange(newResourceFilters)
                              }}
                            />
                          </div>
                        </div>
                      )
                    })
                ) : (
                  <div className='header-filter-tagslist-text'>
                    No filters
                  </div>
                )
              }
            </div>
          ) : null
        }
      </div>
    )
  }
}
