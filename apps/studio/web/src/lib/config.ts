import { useGetConfig } from "#lib/client/config/config";

// Frontend configuration from environment variables
export const getEnvConfig = () => ({
  apiUrl: import.meta.env.VITE_ARCHES_API_HOST || "http://localhost:3001",
  authEnabled: import.meta.env.VITE_ARCHES_AUTH_ENABLED !== "false",
  isDevelopment: import.meta.env.DEV,
  isProduction: import.meta.env.PROD,
  platformUrl:
    import.meta.env.VITE_ARCHES_PLATFORM_URL || "http://localhost:3000",
});

// Hook to get OAuth providers configuration from the backend
export const useOAuthProviders = () => {
  const { data: config, isLoading, error } = useGetConfig();

  const providers = () => {
    if (!config?.data?.auth) return [];

    const availableProviders = [];

    if (config.data?.auth.google?.enabled) {
      availableProviders.push({
        enabled: true,
        id: "google" as const,
        name: "Google",
      });
    }

    if (config.data?.auth.github?.enabled) {
      availableProviders.push({
        enabled: true,
        id: "github" as const,
        name: "GitHub",
      });
    }

    if (config.data?.auth.microsoft?.enabled) {
      availableProviders.push({
        enabled: true,
        id: "microsoft" as const,
        name: "Microsoft",
      });
    }

    return availableProviders;
  };

  return {
    error,
    isLoading,
    providers,
  };
};

// Helper to build OAuth authorization URL
export const buildOAuthUrl = (provider: "google" | "github" | "microsoft") => {
  const { apiUrl } = getEnvConfig();
  const redirectUri = `${window.location.origin}/auth/oauth/callback`;

  // The backend expects the provider in the path and redirect_uri as a query param
  return `${apiUrl}/auth/oauth/${provider}/authorize?redirect_uri=${encodeURIComponent(redirectUri)}`;
};
