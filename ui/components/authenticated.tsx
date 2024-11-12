"use client";

import { useAuth } from "@/hooks/use-auth";
import { useWebsockets } from "@/hooks/useWebsockets";
import { useEffect } from "react";

export function Authenticated({ children }: { children: React.ReactNode }) {
  const { defaultOrgname, getUserFromToken } = useAuth();

  useWebsockets({});

  useEffect(() => {
    if (!defaultOrgname) {
      getUserFromToken();
      return;
    }
  }, [defaultOrgname, getUserFromToken]);

  return <>{children}</>;
}
