import Promise from 'promise'
import navigation from './_navigation'
import summary from './_summary'
import table from './_table'

const mocks = {
  'api/v1/navigation': navigation,
  'api/v1/summary': summary,
  'api/v1/table': table
}

export default function resolveMock (path) {
  return Promise.resolve(mocks[path])
}
