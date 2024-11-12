import type { QueryKey, UseQueryOptions } from "@tanstack/react-query";

import { useAuth } from "@/hooks/use-auth";

import { QueryOperation } from "./archesApiComponents";

export type ArchesApiContext = {
  fetcherOptions: {
    /**
     * Headers to inject in the fetcher
     */
    headers?: {
      authorization?: string;
      "Content-Type"?: string;
    };
    /**
     * Query params to inject in the fetcher
     */
    queryParams?: {};

    signal?: AbortSignal;
  };
  // eslint-disable-next-line no-unused-vars
  queryKeyFn: (operation: QueryOperation) => QueryKey;
  /**
   * Query key manager.
   */
  queryOptions: {
    /**
     * Set this to `false` to disable automatic refetching when the query mounts or changes query keys.
     * Defaults to `true`.
     */
    enabled?: boolean;

    // eslint-disable-next-line no-unused-vars
    onError?: any;
    // eslint-disable-next-line no-unused-vars
    retry?: any;
  };
};

/**
 * Context injected into every react-query hook wrappers
 *
 * @param queryOptions options from the useQuery wrapper
 */
export function useArchesApiContext<
  TQueryFnData = unknown,
  TError = unknown,
  TData = TQueryFnData,
  TQueryKey extends QueryKey = QueryKey,
>(
  _queryOptions?: Omit<
    UseQueryOptions<TQueryFnData, TError, TData, TQueryKey>,
    "queryFn" | "queryKey"
  >
): ArchesApiContext {
  const { getNewRefreshToken, logout, defaultOrgname } = useAuth();
  console.log("RUNNING useArchesApiContext");
  return {
    fetcherOptions: {
      // headers: {
      //
      // },
    },
    queryKeyFn,
    queryOptions: {
      enabled:
        _queryOptions?.enabled !== undefined
          ? !!_queryOptions.enabled
          : !!defaultOrgname,
      retry: async (failureCount: number, error: any) => {
        console.log("RETRYING", failureCount, error);
        if (error?.stack?.statusCode === 401 && failureCount <= 2) {
          await getNewRefreshToken();
          return true;
        } else if (
          (error as any)?.stack?.statusCode === 401 &&
          failureCount > 2
        ) {
          console.log("LOGGING OUT DUE TO TOO MANY RETRIES");
          await logout();
          return false;
        }
        return false;
      },
    },
  };
}

export const queryKeyFn = (operation: QueryOperation) => {
  const queryKey: unknown[] = hasPathParams(operation)
    ? operation.path
        .split("/")
        .filter(Boolean)
        .map((i) => resolvePathParam(i, operation.variables.pathParams))
    : operation.path.split("/").filter(Boolean);

  if (
    operation &&
    operation.variables &&
    operation.variables.queryParams &&
    Object.keys(operation.variables.queryParams).length > 0
  ) {
    queryKey.push(operation.variables.queryParams);
  }

  return queryKey;
};
// Helpers
const resolvePathParam = (key: string, pathParams: Record<string, string>) => {
  if (key.startsWith("{") && key.endsWith("}")) {
    return pathParams[key.slice(1, -1)];
  }
  return key;
};

const hasPathParams = (
  operation: QueryOperation
): operation is {
  variables: { pathParams: Record<string, string> };
} & QueryOperation => {
  return Boolean((operation.variables as any).pathParams);
};

// const hasBody = (
//   operation: QueryOperation,
// ): operation is QueryOperation & {
//   variables: { body: Record<string, unknown> };
// } => {
//   return Boolean((operation.variables as any).body);
// };

// const hasQueryParams = (
//   operation: QueryOperation,
// ): operation is QueryOperation & {
//   variables: { queryParams: Record<string, unknown> };
// } => {
//   return Boolean((operation.variables as any).queryParams);
// };
