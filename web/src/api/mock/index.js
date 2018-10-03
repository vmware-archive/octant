import Promise from 'promise'
import _ from 'lodash'
import navigation from './_navigation'
import summary from './_summary'
import table from './_table'

const mocks = {
  'api/v1/navigation': navigation,
  'api/v1/summary': summary,
  'api/v1/table': table
}

export default function resolveMock (path) {
  if (_.startsWith(path, 'api/v1/content')) return Promise.resolve([summary, table])
  return Promise.resolve(mocks[path])
}
