/* eslint-disable @typescript-eslint/no-empty-object-type */
import { ArchesApiContext } from './archesApiContext'

export const baseUrl = process.env.NEXT_PUBLIC_API_URL

export type ArchesApiFetcherOptions<
  TBody,
  THeaders,
  TQueryParams,
  TPathParams
> = ArchesApiContext['fetcherOptions'] & {
  body?: TBody
  headers?: THeaders
  method: string
  pathParams?: TPathParams
  queryParams?: TQueryParams
  signal?: AbortSignal
  url: string
}

export type ErrorWrapper<TError> =
  | TError
  | { message: string; statusCode: 'unknown' }

export async function archesApiFetch<
  TData,
  TError,
  TBody extends FormData | null | undefined | {},
  THeaders extends {},
  TQueryParams extends {},
  TPathParams extends {}
>({
  body,
  headers,
  method,
  pathParams,
  queryParams,
  signal,
  url
}: ArchesApiFetcherOptions<
  TBody,
  THeaders,
  TQueryParams,
  TPathParams
>): Promise<TData> {
  try {
    const requestHeaders: HeadersInit = {
      'Content-Type': 'application/json',
      ...headers
    }

    /**
     * As the fetch API is being used, when multipart/form-data is specified
     * the Content-Type header must be deleted so that the browser can set
     * the correct boundary.
     * https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects#sending_files_using_a_formdata_object
     */
    if (
      requestHeaders['Content-Type']
        ?.toLowerCase()
        .includes('multipart/form-data')
    ) {
      delete requestHeaders['Content-Type']
    }

    const response = await window.fetch(
      `${baseUrl}${resolveUrl(url, queryParams, pathParams)}`,
      {
        body: body
          ? body instanceof FormData
            ? body
            : JSON.stringify(body)
          : undefined,
        headers: requestHeaders,
        method: method.toUpperCase(),
        signal,
        credentials: 'include'
      }
    )
    if (!response.ok) {
      let error: ErrorWrapper<TError>
      try {
        error = await response.json()
      } catch (e) {
        error = {
          message:
            e instanceof Error
              ? `Unexpected error (${e.message})`
              : 'Unexpected error',
          statusCode: 'unknown' as const
        }
      }

      throw error
    }

    if (response.headers.get('content-type')?.includes('json')) {
      return await response.json()
    } else {
      // if it is not a json response, assume it is a blob and cast it to TData
      return (await response.blob()) as unknown as TData
    }
  } catch (e) {
    const errorObject = e as any
    throw {
      message: errorObject.message,
      statusCode: errorObject.statusCode || 'unknown'
    }
  }
}

const resolveUrl = (
  url: string,
  queryParams: Record<string, string> = {},
  pathParams: Record<string, string> = {}
) => {
  let query = new URLSearchParams(queryParams).toString()
  if (query) query = `?${query}`
  return (
    url.replace(/\{\w*\}/g, (key) => pathParams[key.slice(1, -1)] || '') + query
  )
}
