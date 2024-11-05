"use client";

import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { siteConfig } from "@/config/site";
import { usePathname, useRouter } from "next/navigation";

export const TabsSection = () => {
  const router = useRouter();
  const pathname = usePathname() as string;

  const currentTabs = siteConfig.routes
    .find((route) => pathname.startsWith(route.href))
    ?.children?.filter((tab: any) => tab?.showInTabs);

  if (!currentTabs || currentTabs.length === 0) {
    return <div className="border-b"></div>; // Or return a placeholder if needed
  }

  const activeTab = currentTabs.find((tab) => pathname === tab.href)?.href;

  return (
    <Tabs value={activeTab}>
      <TabsList className="h-8 w-full items-end justify-start rounded-none border-b bg-background">
        {currentTabs.map((tab) => {
          const isActive = tab.href === activeTab;
          console.log(tab.href, activeTab, isActive);
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
  );
};
