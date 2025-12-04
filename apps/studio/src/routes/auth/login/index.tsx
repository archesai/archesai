import type { FormFieldConfig } from "@archesai/ui";
import { Button, GenericForm, Input, toast } from "@archesai/ui";
import { useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useRouter } from "@tanstack/react-router";
import type { JSX } from "react";
import { useState } from "react";
import { TermsIndicator } from "#components/terms-indicator";
import { buildOAuthUrl, getEnvConfig, useOAuthProviders } from "#lib/config";
import type { LoginBody } from "#lib/index";
import { useLogin } from "#lib/index";

export const Route = createFileRoute("/auth/login/")({
  component: LoginPage,
});

function LoginPage(): JSX.Element {
  const router = useRouter();
  const queryClient = useQueryClient();
  const { mutateAsync: login } = useLogin();
  const [authMethod, setAuthMethod] = useState<"password" | "magic-link">(
    "password",
  );
  const [magicLinkSent, setMagicLinkSent] = useState(false);
  const searchParams = Route.useSearch() as {
    message?: string;
    error?: string;
  };

  const { providers: getOAuthProviders } = useOAuthProviders();
  const oauthProviders = getOAuthProviders();

  const handleMagicLinkRequest = async (email: string) => {
    try {
      const { isDevelopment } = getEnvConfig();
      // Use console delivery for development, email for production
      const deliveryMethod = isDevelopment ? "console" : "email";

      const response = await fetch("/api/v1/auth/magic-links/request", {
        body: JSON.stringify({
          deliveryMethod,
          identifier: email,
        }),
        headers: {
          "Content-Type": "application/json",
        },
        method: "POST",
      });

      if (response.ok) {
        setMagicLinkSent(true);
        const message = isDevelopment
          ? "Magic link sent! Check your console for the link."
          : "Magic link sent! Check your email.";
        toast.success(message);
      } else {
        throw new Error("Failed to send magic link");
      }
    } catch (error) {
      console.error("Magic link request failed:", error);
      toast.error("Failed to send magic link. Please try again.");
    }
  };

  const handleOAuthLogin = async (
    provider: "google" | "github" | "microsoft",
  ) => {
    // OAuth flow uses direct navigation to the authorization endpoint
    // The backend will handle the redirect to the OAuth provider
    const authUrl = buildOAuthUrl(provider);

    // Navigate to the OAuth authorization endpoint
    window.location.href = authUrl;
  };

  const formFields: FormFieldConfig<LoginBody>[] =
    authMethod === "password"
      ? [
          {
            defaultValue: "",
            description: "Your email address",
            label: "Email",
            name: "email",
            renderControl: (field) => (
              <Input
                placeholder="Enter your email..."
                {...field}
                type="email"
              />
            ),
          },
          {
            defaultValue: "",
            description: "Your password",
            label: "Password",
            name: "password",
            renderControl: (field) => (
              <Input
                {...field}
                placeholder="Enter your password..."
                type="password"
              />
            ),
          },
        ]
      : [
          {
            defaultValue: "",
            description: "Enter your email to receive a magic link",
            label: "Email",
            name: "email",
            renderControl: (field) => (
              <Input
                placeholder="Enter your email..."
                {...field}
                type="email"
              />
            ),
          },
        ];

  return (
    <>
      {searchParams.message && (
        <div className="mb-4 rounded-md border border-green-200 bg-green-50 p-4 text-green-800 text-sm">
          {searchParams.message}
        </div>
      )}
      {searchParams.error && (
        <div className="mb-4 rounded-md border border-red-200 bg-red-50 p-4 text-red-800 text-sm">
          {searchParams.error}
        </div>
      )}
      {magicLinkSent ? (
        <div className="space-y-4 text-center">
          <div className="rounded-md border border-green-200 bg-green-50 p-4 text-green-800">
            Magic link sent! Check your console (for development) or email.
          </div>
          <Button
            onClick={() => {
              setMagicLinkSent(false);
              setAuthMethod("password");
            }}
            variant="outline"
          >
            Back to login
          </Button>
        </div>
      ) : (
        <GenericForm<LoginBody, never>
          description={
            authMethod === "password"
              ? "Enter your email and password to login"
              : "Enter your email to receive a magic link"
          }
          entityKey="auth"
          fields={formFields}
          isUpdateForm={false}
          onSubmitCreate={async (data) => {
            if (authMethod === "magic-link") {
              await handleMagicLinkRequest(data.email);
              return;
            }

            try {
              await login({
                data: {
                  email: data.email,
                  password: data.password || "",
                },
              });

              // The backend should set httpOnly cookies with the session
              // The response contains access_token and refresh_token
              // The cookie should contain the session ID which will be used by getSession

              // Invalidate the router to trigger a refetch of the session
              await router.invalidate();

              // Clear and refetch queries
              await queryClient.clear();

              // Show success toast
              toast.success("Login successful!");

              // Navigate to home, which will check for the session cookie
              await router.navigate({
                to: "/",
              });
            } catch (error) {
              console.error("Login failed:", error);
              toast.error("Login failed. Please check your credentials.");
              throw error;
            }
          }}
          postContent={
            <div className="space-y-2">
              <div className="text-center text-sm">
                <button
                  className="underline underline-offset-4"
                  onClick={() =>
                    setAuthMethod(
                      authMethod === "password" ? "magic-link" : "password",
                    )
                  }
                  type="button"
                >
                  {authMethod === "password"
                    ? "Sign in with Magic Link"
                    : "Sign in with Password"}
                </button>
              </div>
              <div className="text-center text-sm">
                <Link
                  className="underline underline-offset-4"
                  to="/auth/forgot-password"
                >
                  Forgot your password?
                </Link>
              </div>
              <div className="text-center text-sm">
                Don&apos;t have an account?{" "}
                <Link
                  className="underline underline-offset-4"
                  to="/auth/login"
                >
                  Sign up
                </Link>
              </div>
            </div>
          }
          preContent={
            oauthProviders.length > 0 && (
              <>
                <div className="flex flex-col gap-4">
                  {oauthProviders.find((p) => p.id === "google") && (
                    <Button
                      className="w-full"
                      onClick={() => handleOAuthLogin("google")}
                      type="button"
                      variant="outline"
                    >
                      <svg
                        aria-label="Google"
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M12.48 10.92v3.28h7.84c-.24 1.84-.853 3.187-1.787 4.133-1.147 1.147-2.933 2.4-6.053 2.4-4.827 0-8.6-3.893-8.6-8.72s3.773-8.72 8.6-8.72c2.6 0 4.507 1.027 5.907 2.347l2.307-2.307C18.747 1.44 16.133 0 12.48 0 5.867 0 .307 5.387.307 12s5.56 12 12.173 12c3.573 0 6.267-1.173 8.373-3.36 2.16-2.16 2.84-5.213 2.84-7.667 0-.76-.053-1.467-.173-2.053H12.48z"
                          fill="currentColor"
                        />
                      </svg>
                      Login with Google
                    </Button>
                  )}
                  {oauthProviders.find((p) => p.id === "github") && (
                    <Button
                      className="w-full"
                      onClick={() => handleOAuthLogin("github")}
                      type="button"
                      variant="outline"
                    >
                      <svg
                        aria-label="GitHub"
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"
                          fill="currentColor"
                        />
                      </svg>
                      Login with GitHub
                    </Button>
                  )}
                  {oauthProviders.find((p) => p.id === "microsoft") && (
                    <Button
                      className="w-full"
                      onClick={() => handleOAuthLogin("microsoft")}
                      type="button"
                      variant="outline"
                    >
                      <svg
                        aria-label="Microsoft"
                        role="img"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M11.4 24H0V12.6h11.4V24zM24 24H12.6V12.6H24V24zM11.4 11.4H0V0h11.4v11.4zm12.6 0H12.6V0H24v11.4z"
                          fill="currentColor"
                        />
                      </svg>
                      Login with Microsoft
                    </Button>
                  )}
                </div>
                <div className="relative text-center text-sm after:absolute after:inset-0 after:top-1/2 after:z-0 after:flex after:items-center after:border-border after:border-t">
                  <span className="relative z-10 bg-card px-2 text-muted-foreground">
                    Or continue with
                  </span>
                </div>
              </>
            )
          }
          showCard={true}
          title="Login"
        />
      )}
      <TermsIndicator />
    </>
  );
}
