import { ArchesLogo } from "@archesai/ui";

import {
  createFileRoute,
  Link,
  Outlet,
  redirect,
} from "@tanstack/react-router";
import type { JSX } from "react";
import { getEnvConfig } from "#lib/config";

export const Route = createFileRoute("/auth")({
  beforeLoad: ({ context }) => {
    const REDIRECT_URL = "/";
    const { authEnabled } = getEnvConfig();
    // Redirect away from auth pages if auth is disabled or user is already logged in
    if (!authEnabled || context.session?.data) {
      throw redirect({
        to: REDIRECT_URL,
      });
    }
    return {
      redirectUrl: REDIRECT_URL,
    };
  },
  component: AuthenticationLayout,
});

function AuthenticationLayout(): JSX.Element {
  return (
    <div className="flex h-svh flex-col items-center justify-center gap-4">
      <Link to="/">
        <ArchesLogo scale={1.5} />
      </Link>
      <Outlet />
    </div>
  );
}
