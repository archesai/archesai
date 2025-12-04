/// <reference types="vite/client" />

const getBody = <T>(c: Request | Response): Promise<T> => {
  const contentType = c.headers.get("content-type");

  if (contentType?.includes("application/json")) {
    return c.json() satisfies Promise<T>;
  }

  if (contentType?.includes("application/pdf")) {
    return c.blob() as Promise<T>;
  }

  return c.text() as Promise<T>;
};

const getUrl = (contextUrl: string): string => {
  // Use relative URL for SSR compatibility (Vite proxy handles routing)
  if (typeof window === "undefined") {
    // Server-side: use configured host
    const host = import.meta.env.VITE_ARCHES_API_HOST;
    if (!host) {
      throw new Error("host URL is not configured.");
    }
    try {
      const requestUrl = new URL(`http://${host}/api/v1${contextUrl}`);
      return requestUrl.toString();
    } catch (error) {
      throw new Error(`could not parse url: ${error}`);
    }
  }

  // Client-side: use Vite proxy (relative URL)
  return `/api/v1${contextUrl}`;
};

import type { Problem } from "./client/orval.schemas";

const isProblem = (obj: unknown): obj is Problem => {
  return (
    obj !== null &&
    typeof obj === "object" &&
    "type" in obj &&
    typeof obj.type === "string"
  );
};

export const customFetch = async <T>(
  url: string,
  options: RequestInit,
): Promise<T> => {
  const requestUrl = getUrl(url);

  const requestInit: RequestInit = {
    ...options,
    credentials: "include",
    headers: new Headers(options.headers),
  };

  try {
    const response = await fetch(requestUrl, requestInit);

    if (!response.ok) {
      let problem: Problem;
      try {
        const errorData = await getBody<unknown>(response.clone());
        if (isProblem(errorData)) {
          problem = errorData;
        } else {
          throw new Error("Invalid problem format");
        }
      } catch {
        // If we can't parse as Problem schema, create a basic one
        problem = {
          detail: `HTTP error ${response.status}`,
          status: response.status,
          title: response.statusText || "Unknown Error",
          type: "about:blank",
        };
      }

      throw problem;
    }

    const data = await getBody<T>(response);
    return data as T;
  } catch (error) {
    console.error("Fetch error:", error);
    // If it's already a Problem, re-throw it
    if (isProblem(error)) {
      throw error;
    }

    // Check if it's a Node.js network error (ECONNREFUSED, etc.)
    const isNetworkError =
      error instanceof Error &&
      ("code" in error ||
        "errno" in error ||
        "syscall" in error ||
        error.message.includes("ECONNREFUSED") ||
        error.message.includes("fetch failed"));

    // Handle network errors, timeouts, etc. as Problem format
    const networkProblem: Problem = {
      detail: isNetworkError
        ? `Connection failed: ${error instanceof Error ? error.message : "Unknown network error"}`
        : error instanceof Error
          ? error.message
          : "Unknown network error",
      status: isNetworkError ? 503 : 500,
      title: isNetworkError
        ? "Service Unavailable"
        : error instanceof Error
          ? error.name
          : "Network Error",
      type: "about:blank",
    };

    throw networkProblem;
  }
};
