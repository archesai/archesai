import { toast } from "@archesai/ui";
import {
  keepPreviousData,
  MutationCache,
  notifyManager,
  QueryClient,
} from "@tanstack/react-query";
import { createRouter as createTanStackRouter } from "@tanstack/react-router";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";

import { DefaultCatchBoundary } from "#components/default-catch-boundary";
import NotFound from "#components/not-found";
import { routeTree } from "#routeTree.gen";

export function createRouter() {
  if (typeof document !== "undefined") {
    notifyManager.setScheduler(window.requestAnimationFrame);
  }

  const queryClient: QueryClient = new QueryClient({
    defaultOptions: {
      queries: {
        placeholderData: keepPreviousData,
        refetchOnReconnect: () => !queryClient.isMutating(),
        refetchOnWindowFocus: false,
        retry: false,
        staleTime: 1000 * 60 * 2, // 2 minutes,
      },
    },
    mutationCache: new MutationCache({
      onError: (error) => {
        toast.error("An error occurred", {
          className: "bg-red-500 text-white",
          description: error.message.replaceAll(":", ""),
        });
      },
      onSettled: () => {
        if (queryClient.isMutating() === 1) {
          return queryClient.invalidateQueries();
        } else {
          return;
        }
      },
      onSuccess: () => {
        toast.success("Success", {
          className: "bg-green-500 text-white",
        });
      },
    }),
  });

  const router = createTanStackRouter({
    context: {
      queryClient,
      session: null,
    },
    defaultErrorComponent: DefaultCatchBoundary,
    defaultNotFoundComponent: () => <NotFound />,
    defaultPendingComponent: () => {
      return (
        <div className="flex h-full w-full items-center justify-center">
          <div className="h-12 w-12 animate-spin rounded-full border-4 border-primary border-t-transparent text-primary" />
        </div>
      );
    },
    defaultPreload: "intent",
    // https://tanstack.com/router/latest/docs/framework/react/guide/data-loading#passing-all-loader-events-to-an-external-cache
    defaultPreloadStaleTime: 0,
    defaultStructuralSharing: true,
    routeTree,
    scrollRestoration: true,
  });

  setupRouterSsrQueryIntegration({
    queryClient,
    router,
  });

  return router;
}

declare module "@tanstack/react-router" {
  interface Register {
    router: ReturnType<typeof createRouter>;
  }
}
