// NOTE: Supports cases where `content-type` is other than `json`
const getBody = <T>(c: Request | Response): Promise<T> => {
  const contentType = c.headers.get('content-type')

  if (contentType?.includes('application/json')) {
    return c.json() as Promise<T>
  }

  if (contentType?.includes('application/pdf')) {
    return c.blob() as Promise<T>
  }

  return c.text() as Promise<T>
}

// NOTE: Update just base url
const getUrl = (contextUrl: string): string => {
  const baseUrl = 'https://api.archesai.dev'
  const requestUrl = new URL(`${baseUrl}${contextUrl}`)

  return requestUrl.toString()
}

// NOTE: Add headers
const getHeaders = (headers?: HeadersInit): HeadersInit => {
  return {
    ...headers
    // Authorization: 'token',
    // 'Content-Type': 'multipart/form-data'
  }
}

export const customFetch = async <T>(
  url: string,
  options: RequestInit
): Promise<T> => {
  const requestUrl = getUrl(url)
  const requestHeaders = getHeaders(options.headers)

  const requestInit: RequestInit = {
    ...options,
    credentials: 'include', // Ensure credentials are included
    headers: requestHeaders
  }

  const response = await fetch(requestUrl, requestInit)
  const data = await getBody<T>(response)
  if (!response.ok) {
    throw new Error(
      `Request failed with status ${response.status.toString()}: ${response.statusText}`
    )
  }

  return data as T
}
