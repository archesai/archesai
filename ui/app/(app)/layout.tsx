"use client";

import { PageHeader } from "@/components/page-header";
import { AppSidebar } from "@/components/sidebar/app-sidebar";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
import { useAuth } from "@/hooks/useAuth";
import { useWebsockets } from "@/hooks/useWebsockets";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const router = useRouter();
  const { accessToken, getUserFromToken, user } = useAuth();
  const [isHydrated, setIsHydrated] = useState(false);

  useWebsockets({});

  useEffect(() => {
    setIsHydrated(true);
  }, [isHydrated]);

  useEffect(() => {
    if (!isHydrated) return;
    if (!accessToken) {
      router.push("/login");
      return;
    }
    if (!user?.defaultOrgname) {
      getUserFromToken();
      return;
    }
  }, [accessToken, user, getUserFromToken, router, isHydrated]);

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <main className="flex max-h-screen flex-1 flex-col bg-gray-50 dark:bg-neutral-950">
          <PageHeader />
          <div className="flex-1 overflow-auto p-4">{children}</div>
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
