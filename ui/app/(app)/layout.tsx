"use client";

import { VerifyEmailAlert } from "@/components/email-verify";
import { Sidebar } from "@/components/sidebar";
import { TabsSection } from "@/components/tabs-section";
import { useAuth } from "@/hooks/useAuth";
import { useFullScreen } from "@/hooks/useFullScreen";
import { useSidebar } from "@/hooks/useSidebar";
import { useWebsockets } from "@/hooks/useWebsockets";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { PageHeader } from "../../components/page-header";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const router = useRouter();
  const { accessToken, getUserFromToken, user } = useAuth();
  const { isCollapsed } = useSidebar();
  const { isFullScreen } = useFullScreen();
  const [isHydrated, setIsHydrated] = useState(false);

  useWebsockets({});

  useEffect(() => {
    if (isHydrated) return;
    setIsHydrated(true);
  }, [isHydrated]);

  useEffect(() => {
    if (!isHydrated) return;
    if (!accessToken) {
      router.push("/auth/login");
      return;
    }
    if (!user?.defaultOrgname) {
      getUserFromToken();
      return;
    }
  }, [accessToken, user, getUserFromToken, router, isHydrated]);

  return (
    <div
      className={`flex h-screen w-full flex-col transition-all duration-200 md:grid ${
        !isCollapsed ? "md:grid-cols-[250px_1fr]" : "md:grid-cols-[65px_1fr]"
      }`}
    >
      {/* Sidebar for desktop */}
      <div className="hidden border-r shadow-sm md:block">
        <Sidebar />
      </div>

      {/* Main Content */}
      <main className="flex max-h-screen flex-1 flex-col bg-gray-50 dark:bg-black">
        {user && !user.emailVerified && <VerifyEmailAlert />}
        {!isFullScreen && (
          <>
            <PageHeader />
            <TabsSection />
          </>
        )}
        <div className="flex-1 overflow-auto p-4">{children}</div>
      </main>
    </div>
  );
}
