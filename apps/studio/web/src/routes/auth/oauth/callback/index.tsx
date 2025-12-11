import { Button } from "@archesai/ui";
import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import type { JSX } from "react";
import { useEffect } from "react";
import { useOauthCallback } from "#lib/index";

type OAuthProvider = "google" | "github" | "microsoft";

export const Route = createFileRoute("/auth/oauth/callback/")({
  component: OAuthCallbackPage,
});

function OAuthCallbackPage(): JSX.Element {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const searchParams = Route.useSearch() as Record<string, string | undefined>;

  const provider = searchParams?.provider as OAuthProvider | undefined;
  const code = searchParams?.code;
  const state = searchParams?.state;
  const error = searchParams?.error;
  const errorDescription = searchParams?.error_description;

  const isValidProvider =
    provider && ["google", "github", "microsoft"].includes(provider);

  const {
    data,
    isLoading,
    isError,
    error: callbackError,
  } = useOauthCallback(
    provider as OAuthProvider,
    isValidProvider && (code || error)
      ? {
          ...(code && { code }),
          ...(error && { error }),
          ...(errorDescription && { error_description: errorDescription }),
          ...(state && { state }),
        }
      : undefined,
  );

  useEffect(() => {
    if (!isValidProvider) {
      navigate({
        search: { error: "Invalid authentication provider" },
        to: "/auth/login",
      });
    }
  }, [isValidProvider, navigate]);

  useEffect(() => {
    if (data) {
      queryClient.invalidateQueries();
      navigate({ to: "/" });
    }
  }, [data, navigate, queryClient]);

  useEffect(() => {
    if (error || callbackError) {
      const errorMsg =
        errorDescription || callbackError?.detail || "Authentication failed";
      navigate({
        search: { error: errorMsg },
        to: "/auth/login",
      });
    }
  }, [error, errorDescription, callbackError, navigate]);

  if (!isValidProvider) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h2 className="font-semibold text-destructive text-lg">
            Invalid Provider
          </h2>
          <p className="mt-2 text-muted-foreground text-sm">
            The authentication provider is invalid or missing.
          </p>
        </div>
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h2 className="font-semibold text-lg">Authenticating...</h2>
          <p className="mt-2 text-muted-foreground text-sm">
            Please wait while we complete your {provider} login.
          </p>
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <h2 className="font-semibold text-destructive text-lg">
            Authentication Failed
          </h2>
          <p className="mt-2 text-muted-foreground text-sm">
            {callbackError?.detail ||
              "An error occurred during authentication."}
          </p>
          <Button
            className="mt-4 rounded-md bg-primary px-4 py-2 font-medium text-primary-foreground text-sm hover:bg-primary/90"
            onClick={() => navigate({ to: "/auth/login" })}
          >
            Return to Login
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <h2 className="font-semibold text-lg">Processing...</h2>
      </div>
    </div>
  );
}
