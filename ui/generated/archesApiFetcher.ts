import { ArchesApiContext } from "./archesApiContext";

export const baseUrl = process.env.NEXT_PUBLIC_API_URL;

export type ErrorWrapper<TError> =
  | { payload: string; stack: any; status: "unknown" }
  | TError;

export type ArchesApiFetcherOptions<
  TBody,
  THeaders,
  TQueryParams,
  TPathParams,
> = {
  body?: TBody;
  headers?: THeaders;
  method: string;
  pathParams?: TPathParams;
  queryParams?: TQueryParams;
  signal?: AbortSignal;
  url: string;
} & ArchesApiContext["fetcherOptions"];

export async function archesApiFetch<
  TData,
  TError,
  TBody extends {} | FormData | null | undefined,
  THeaders extends {},
  TQueryParams extends {},
  TPathParams extends {},
>({
  body,
  headers,
  method,
  pathParams,
  queryParams,
  signal,
  url,
}: ArchesApiFetcherOptions<
  TBody,
  THeaders,
  TQueryParams,
  TPathParams
>): Promise<TData> {
  try {
    const requestHeaders: HeadersInit = {
      "Content-Type": "application/json",
      ...headers,
    };

    /**
     * As the fetch API is being used, when multipart/form-data is specified
     * the Content-Type header must be deleted so that the browser can set
     * the correct boundary.
     * https://developer.mozilla.org/en-US/docs/Web/API/FormData/Using_FormData_Objects#sending_files_using_a_formdata_object
     */
    if (
      requestHeaders["Content-Type"]
        .toLowerCase()
        .includes("multipart/form-data")
    ) {
      delete requestHeaders["Content-Type"];
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
      }
    );
    if (!response.ok) {
      let error: ErrorWrapper<TError>;
      try {
        error = await response.json();
      } catch (e) {
        error = {
          payload:
            e instanceof Error
              ? `Unexpected error (${e.message})`
              : "Unexpected error",
          stack: {},
          status: "unknown" as const,
        };
      }

      throw error;
    }

    if (response.headers.get("content-type")?.includes("json")) {
      return await response.json();
    } else {
      // if it is not a json response, assume it is a blob and cast it to TData
      return (await response.blob()) as unknown as TData;
    }
  } catch (e) {
    let err = e as any;
    const errorObject: Error = {
      message:
        e instanceof Error ? `Network error (${err.message})` : err.message,
      name: "unknown" as const,
      stack: e as string,
    };
    throw errorObject;
  }
}

const resolveUrl = (
  url: string,
  queryParams: Record<string, string> = {},
  pathParams: Record<string, string> = {}
) => {
  let query = new URLSearchParams(queryParams).toString();
  if (query) query = `?${query}`;
  return url.replace(/\{\w*\}/g, (key) => pathParams[key.slice(1, -1)]) + query;
};
