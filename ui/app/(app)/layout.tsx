// app/layout.tsx
"use client";

import { CommandMenu } from "@/components/command-menu";
import { VerifyEmailAlert } from "@/components/email-verify";
import { ModeToggle } from "@/components/mode-toggle";
import { Sidebar } from "@/components/sidebar";
import { TabsSection } from "@/components/tabs-section";
import { Button } from "@/components/ui/button";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { UserButton } from "@/components/user-button";
import { useAuth } from "@/hooks/useAuth";
import { useFullScreen } from "@/hooks/useFullScreen";
import { useSidebar } from "@/hooks/useSidebar";
import { useWebsockets } from "@/hooks/useWebsockets";
import { Menu } from "lucide-react";
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
      className={`flex flex-col md:grid h-screen w-full transition-all duration-200 ${
        !isCollapsed ? "md:grid-cols-[250px_1fr]" : "md:grid-cols-[65px_1fr]"
      }`}
    >
      {/* Sidebar for desktop */}
      <div className="hidden border-r shadow-sm md:block">
        <Sidebar />
      </div>

      {/* Sidebar for mobile */}
      <div className="flex md:hidden min-h-16 items-center px-6 z-50">
        <Sheet>
          <div className="flex items-center justify-between w-full">
            <SheetTrigger asChild>
              <Button
                className="shrink-0 md:hidden"
                size="icon"
                variant="outline"
              >
                <Menu className="h-5 w-5" />
                <span className="sr-only">Toggle navigation menu</span>
              </Button>
            </SheetTrigger>
            <div className="flex items-center gap-3">
              <CommandMenu />
              <ModeToggle />
              <UserButton size="sm" />
            </div>
          </div>

          <SheetContent className="w-64 p-0" side="left">
            <Sidebar />
          </SheetContent>
        </Sheet>
      </div>

      {/* Main Content */}
      <div className="flex flex-col flex-1 bg-gray-50 dark:bg-black">
        <main className="flex flex-1 flex-col overflow-hidden max-h-screen">
          {user && !user.emailVerified && <VerifyEmailAlert />}
          {!isFullScreen && (
            <>
              <PageHeader description={""} title={""} />
              <TabsSection />
            </>
          )}
          {/* Updated to allow scrolling when content overflows */}
          <div className="px-6 md:py-6 py-4 flex-1 overflow-auto">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}
