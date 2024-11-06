"use client";

import { VerifyEmailAlert } from "@/components/email-verify";
import { AppSidebar } from "@/components/sidebar/app-sidebar";
import { TabsSection } from "@/components/tabs-section";
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar";
// import { siteConfig } from "@/config/site";
import { useAuth } from "@/hooks/useAuth";
import { useWebsockets } from "@/hooks/useWebsockets";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { PageHeader } from "../../components/page-header";

export default function AppLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  // const pathname = usePathname();
  // // combine all the routes from siteConfig
  // const routes = siteConfig.routes
  //   .map((route) => [route, ...(route.children || [])])
  //   .flat();
  // // find the current route
  // const currentRoute = routes.find((route) => pathname === route.href);
  // // get the title and description from the current route
  // const title = currentRoute?.title;
  // const description = currentRoute?.description;

  const router = useRouter();
  const { accessToken, getUserFromToken, user } = useAuth();
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
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <main className="flex max-h-screen flex-1 flex-col bg-gray-50 dark:bg-neutral-950">
          {user && !user.emailVerified && <VerifyEmailAlert />}
          <PageHeader />
          <TabsSection />
          <div className="flex-1 overflow-auto p-4">
            {/* {description && (
              <div className="flex h-16 flex-col">
                <div className="text-xl font-semibold text-foreground">
                  {title}
                </div>
                <div className="text-sm text-muted-foreground">
                  {description}
                </div>
              </div>
            )} */}
            {children}
          </div>
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
