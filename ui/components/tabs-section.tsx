"use client";

import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { siteConfig } from "@/config/site";
import { usePathname, useRouter } from "next/navigation";

export const TabsSection = () => {
  const router = useRouter();
  const pathname = usePathname();

  // Determine which tabs to display based on the current route pattern
  let currentTabs: {
    href: string;
    Icon: any;
    tab?: string;
    title: string;
  }[] = [];

  // Iterate over siteConfig.links to find the matching route pattern
  const entries = Object.entries(siteConfig.links);
  // sort the entries by the length of the route pattern in descending order
  entries.sort(([a], [b]) => b.length - a.length);
  for (const [routePattern, links] of entries) {
    if (
      pathname.startsWith(routePattern) ||
      getOriginalPath(pathname).startsWith(routePattern)
    ) {
      // Filter links that are intended for tabs (links without 'section' property)
      currentTabs = links.filter((link) => link.tab);
      break;
    }
  }

  // If no tabs are found, return null or a placeholder
  if (currentTabs.length === 0) {
    return <div className="border-b shadow-sm"></div>; // Or return a placeholder if needed
  }

  // Determine the active tab
  const activeTab = currentTabs.find((tab) => {
    const regex = convertPatternToRegex(tab.href);
    return regex.test(pathname) || getOriginalPath(pathname) === tab.href;
  })?.href;

  return (
    <Tabs className="relative w-full bg-background" value={activeTab}>
      <TabsList className="w-full justify-start rounded-none border-b bg-transparent p-0">
        {currentTabs.map((tab) => {
          const isActive = tab.href === activeTab;
          return (
            <TabsTrigger
              className={`relative rounded-none border-b-2 px-4 pb-3 pt-2 font-semibold shadow-none transition-all 
                ${
                  isActive
                    ? "border-b-primary text-foreground"
                    : "text-muted-foreground"
                }
              `}
              key={tab.href}
              onClick={() =>
                router.push(replaceDynamicSegments(tab.href, pathname))
              }
              value={tab.href}
            >
              {tab.tab || tab.title}
            </TabsTrigger>
          );
        })}
      </TabsList>
    </Tabs>
  );
};

// Helper function to replace dynamic segments in the href with actual values from the pathname
const replaceDynamicSegments = (href: string, pathname: string): string => {
  const hrefParts = href.split("/");
  const pathnameParts = pathname.split("/");

  return hrefParts
    .map((part, index) => {
      if (part.startsWith("[") && part.endsWith("]")) {
        return pathnameParts[index] || part;
      }
      return part;
    })
    .join("/");
};

const convertPatternToRegex = (pattern: string): RegExp => {
  // Replace dynamic segments like [param] with a regex pattern
  const patternWithPlaceholders = pattern.replace(/\[([^\]]+)\]/g, "[^/]+");

  // Escape special regex characters, except for the regex components we've just added
  const escapedPattern = patternWithPlaceholders.replace(
    /[-\\^$*+?.()|{}]/g,
    "\\$&"
  );

  return new RegExp(`^${escapedPattern}$`);
};

export function getOriginalPath(pathname: string) {
  // Define patterns for typical dynamic segments
  const patterns = [
    {
      regex:
        /^\/chatbots\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/,
      replacement: "/chatbots/[chatbotId]",
    },
    {
      regex:
        /^\/users\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}/,
      replacement: "/users/[userId]",
    },
    {
      regex:
        /^\/chatbots\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\/threads/,
      replacement: "/chatbots/[chatbotId]/threads",
    },
    {
      regex:
        /^\/chatbots\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\/chat/,
      replacement: "/chatbots/[chatbotId]/chat",
    },
    {
      regex:
        /^\/chatbots\/[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}\/configuration/,
      replacement: "/chatbots/[chatbotId]/configuration",
    },
  ];
  // Check each pattern to see if there's a match
  for (const pattern of patterns) {
    if (pathname.match(pattern.regex)) {
      return pathname.replace(pattern.regex, pattern.replacement);
    }
  }

  // Return the original pathname if no patterns match
  return pathname;
}
