import { createIsomorphicFn } from '@tanstack/react-start'
import { getWebRequest } from '@tanstack/react-start/server'

export const getIsomorphicHeaders = createIsomorphicFn()
  .client(() => ({}))
  .server(() => getWebRequest().headers)
