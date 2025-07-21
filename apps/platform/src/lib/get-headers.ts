import { createIsomorphicFn, createServerFn } from '@tanstack/react-start'
import { getHeaders as getServerHeaders } from '@tanstack/react-start/server'

export const getHeaders = createServerFn({ method: 'GET' }).handler(() => {
  return getServerHeaders()
})

export const getHeadersIsomorphic = createIsomorphicFn()
  .client(() => {
    return getHeaders()
  })
  .server(() => {
    return getServerHeaders()
  })
