import { getContents } from 'api'

export default async function (pathname, namespace) {
  try {
    const payload = await getContents(pathname, namespace)
    if (payload) {
      return {
        contents: payload.contents,
        title: payload.title,
        loading: false,
        error: false
      }
    }
  } catch (e) {
    return { loading: false, error: true }
  }
}
