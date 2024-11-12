"use client";

import { useSidebar } from "@/components/ui/sidebar";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { siteConfig } from "@/config/site";
import { useAuth } from "@/hooks/use-auth";
import { Menu } from "lucide-react";
import { usePathname, useRouter } from "next/navigation";

import { Button } from "../../ui/button";
import { CommandMenu } from "./command-menu";
import { VerifyEmailAlert } from "./email-verify";
import { ModeToggle } from "./mode-toggle";
import { UserButton } from "./user-button";

export const PageHeader = () => {
  const { toggleSidebar } = useSidebar();
  const router = useRouter();
  const pathname = usePathname() as string;
  const { user } = useAuth();

  // combine all the routes from siteConfig
  const routes = siteConfig.routes
    .map((route) => [route, ...(route.children || [])])
    .flat();

  // find the current route
  const currentRoute = routes.find((route) => pathname === route.href);
  // get the title and description from the current route
  const title = currentRoute?.title;
  const description = currentRoute?.description;

  const currentTabs = siteConfig.routes
    .find((route) => pathname.startsWith(route.href))
    ?.children?.filter((tab: any) => tab?.showInTabs);
  const activeTab = currentTabs?.find((tab) => pathname === tab.href)?.href;

  return (
    <>
      {user && !user.emailVerified && <VerifyEmailAlert />}

      <header className="flex w-full items-center justify-between bg-background p-3 py-3">
        <Button
          className="mr-3 flex h-8 w-8"
          onClick={toggleSidebar}
          size="icon"
          variant="secondary"
        >
          <Menu className="h-5 w-5" />
        </Button>
        <div className="flex flex-1 items-center justify-end gap-3">
          <CommandMenu />
          <ModeToggle />
          <UserButton size="sm" />
        </div>
      </header>

      {!currentTabs || currentTabs.length === 0 ? (
        <div className="border-b" />
      ) : (
        <Tabs value={activeTab}>
          <TabsList className="h-8 w-full items-end justify-start rounded-none border-b bg-background">
            {currentTabs.map((tab) => {
              const isActive = tab.href === activeTab;
              return (
                <TabsTrigger
                  className={`relative h-8 font-normal shadow-none transition-all hover:bg-muted [&::after]:absolute [&::after]:bottom-0 [&::after]:left-0 [&::after]:h-0.5 [&::after]:bg-primary [&::after]:transition-all [&::after]:content-[''] ${isActive ? "text-foreground [&::after]:w-full" : "text-muted-foreground [&::after]:w-0"}`}
                  key={tab.href}
                  onClick={() => {
                    router.push(tab.href);
                  }}
                  value={tab.href}
                >
                  {tab.title}
                </TabsTrigger>
              );
            })}
          </TabsList>
        </Tabs>
      )}

      <div className="px-4 pt-4">
        <p className="text-xl font-semibold text-foreground">{title}</p>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>
    </>
  );
};
