import type { QueryKey, UseQueryOptions } from '@tanstack/react-query'

import { useAuth } from '@/hooks/use-auth'

import { QueryOperation } from './archesApiComponents'

export type ArchesApiContext = {
  fetcherOptions: {
    /**
     * Headers to inject in the fetcher
     */
    headers?: {
      authorization?: string
      'Content-Type'?: string
    }
    /**
     * Query params to inject in the fetcher
     */
    queryParams?: object
  }
  /**
   * Query key manager.
   */
  queryKeyFn: (operation: QueryOperation) => QueryKey
  queryOptions: {
    /**
     * Set this to `false` to disable automatic refetching when the query mounts or changes query keys.
     * Defaults to `true`.
     */
    enabled?: boolean

    onError?: any

    retry?: any
  }
}

/**
 * Context injected into every react-query hook wrappers
 *
 * @param queryOptions options from the useQuery wrapper
 */
export function useArchesApiContext<
  TQueryFnData = unknown,
  TError = unknown,
  TData = TQueryFnData,
  TQueryKey extends QueryKey = QueryKey
>(
  _queryOptions?: Omit<
    UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>,
    'queryFn' | 'queryKey'
  >
): ArchesApiContext {
  const { defaultOrgname, logout, setStatus } = useAuth()
  return {
    fetcherOptions: {},
    queryKeyFn,
    queryOptions: {
      enabled:
        _queryOptions?.enabled !== undefined
          ? !!_queryOptions.enabled
          : !!defaultOrgname,
      retry: async (failureCount: number, error: any) => {
        console.log('FAILED, RETRYING', failureCount, error)
        if (error?.statusCode === 401 && failureCount <= 2) {
          setStatus('Unauthenticated')
          return true
        } else if (error?.statusCode === 401 && failureCount > 2) {
          console.log('Too many retries, logging out')
          await logout()
          return false
        }
        return false
      }
    }
  }
}

export const queryKeyFn = (operation: QueryOperation) => {
  const queryKey: unknown[] = hasPathParams(operation)
    ? operation.path
        .split('/')
        .filter(Boolean)
        .map((i) => resolvePathParam(i, operation.variables.pathParams))
    : operation.path.split('/').filter(Boolean)

  if (
    operation.variables.queryParams &&
    Object.keys(operation.variables.queryParams).length > 0
  ) {
    queryKey.push(operation.variables.queryParams)
  }

  return queryKey
}
// Helpers
const resolvePathParam = (key: string, pathParams: Record<string, string>) => {
  if (key.startsWith('{') && key.endsWith('}')) {
    return pathParams[key.slice(1, -1)]
  }
  return key
}

const hasPathParams = (
  operation: QueryOperation
): operation is QueryOperation & {
  variables: { pathParams: Record<string, string> }
} => {
  return Boolean((operation.variables as any).pathParams)
}

// const hasBody = (
//   operation: QueryOperation
// ): operation is QueryOperation & {
//   variables: { body: Record<string, unknown> }
// } => {
//   return Boolean((operation.variables as any).body)
// }

// const hasQueryParams = (
//   operation: QueryOperation
// ): operation is QueryOperation & {
//   variables: { queryParams: Record<string, unknown> }
// } => {
//   return Boolean((operation.variables as any).queryParams)
// }
