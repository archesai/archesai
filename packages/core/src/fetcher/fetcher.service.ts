export interface FetchOptions {
  body?: Record<string, unknown>
  headers?: Record<string, string>
  method?: 'DELETE' | 'GET' | 'HEAD' | 'POST' | 'PUT'
}

/**
 * A service to fetch data from a remote server
 */
export class FetcherService {
  public delete<T>(url: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(url, {
      ...(headers ?? {}),
      method: 'DELETE'
    })
  }

  public get<T>(url: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(url, { ...(headers ?? {}), method: 'GET' })
  }

  public head<T>(url: string, headers?: Record<string, string>): Promise<T> {
    return this.request<T>(url, { ...(headers ?? {}), method: 'HEAD' })
  }

  public post<T>(
    url: string,
    body: Record<string, unknown>,
    headers?: Record<string, string>
  ): Promise<T> {
    return this.request<T>(url, { body, ...(headers ?? {}), method: 'POST' })
  }

  public put<T>(
    url: string,
    body: Record<string, unknown>,
    headers?: Record<string, string>
  ): Promise<T> {
    return this.request<T>(url, { body, ...(headers ?? {}), method: 'PUT' })
  }

  private async request<T>(
    url: string,
    options: FetchOptions = {}
  ): Promise<T> {
    const response: Response = await fetch(url, {
      ...(options.body ? { body: JSON.stringify(options.body) } : {}),
      headers: {
        'Content-Type': 'application/json',
        ...(options.headers ?? {})
      },
      method: options.method ?? 'GET'
    })

    if (!response.ok) {
      throw new Error(
        `HTTP Error: ${response.status.toString()} - ${response.statusText}`
      )
    }

    return response.json() as Promise<T>
  }
}
