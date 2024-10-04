"use client";

import { LogoSVG } from "@/components/logo-svg";
import { Button } from "@/components/ui/button";
import { UserButton } from "@/components/user-button";
import { siteConfig } from "@/config/site";
import { useSidebar } from "@/hooks/useSidebar";
import { Menu } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { useTheme } from "next-themes";

import { CreditQuota } from "./credit-quota";

export const Sidebar = () => {
  const { isCollapsed, toggleSidebar } = useSidebar();
  const { resolvedTheme } = useTheme();
  const pathname = usePathname();

  // Sections to display in the sidebar
  const sidebarSections = ["home", "create", "data", "settings"];

  // Collect links grouped by section
  const linksBySection: {
    [key: string]: {
      href: string;
      Icon: any;
      title: string;
    }[];
  } = {};

  // Initialize linksBySection
  sidebarSections.forEach((section) => {
    linksBySection[section] = [];
  });

  // Iterate over siteConfig.links and collect links by section
  Object.values(siteConfig.links).forEach((links) => {
    links.forEach((link) => {
      if (link.section && sidebarSections.includes(link.section)) {
        linksBySection[link.section].push(link);
      }
    });
  });

  return (
    <nav className="items-start px-2.5 pb-2.5 pt-1.5 text-sm font-md h-full">
      <div className="flex flex-col justify-between max-h-screen h-full opacity-100">
        <div>
          <div
            className={`flex items-center py-3 ${
              isCollapsed ? "justify-center" : "justify-between"
            }`}
          >
            {!isCollapsed && (
              <LogoSVG
                fill={resolvedTheme === "dark" ? "#FFFFFF" : "#000"}
                scale={0.08}
                size={"lg"}
              />
            )}

            <Button onClick={toggleSidebar} size="icon" variant="secondary">
              <Menu className="h-5 w-5" />
            </Button>
          </div>

          {/* Render specific sidebar sections */}
          {sidebarSections.map((section) => (
            <div className="mt-4" key={section}>
              {section !== "home" && (
                <h2
                  className={`text-xs px-1 uppercase font-semibold text-muted-foreground mb-2 whitespace-nowrap inter opacity-100 ${
                    isCollapsed ? "text-white" : ""
                  }`}
                >
                  {section.charAt(0).toUpperCase() + section.slice(1)}
                </h2>
              )}

              {linksBySection[section].map(({ href, Icon, title }) => (
                <Link
                  className={`${
                    isCollapsed ? "justify-center" : ""
                  } flex items-center text-md font-medium gap-3 rounded-lg px-2 py-2.5 my-1 ${
                    pathname === href || pathname.startsWith(href)
                      ? "bg-muted"
                      : ""
                  } transition-all hover:bg-muted relative group`}
                  href={href}
                  key={href}
                >
                  <Icon className="h-5 w-5" strokeWidth={1.5} />
                  <span
                    className={`${isCollapsed ? "hidden" : "block"} mt-0.5`}
                  >
                    {title}
                  </span>
                  {isCollapsed && (
                    <span className="absolute left-full ml-2 whitespace-nowrap bg-gray-800 text-white text-xs rounded-md px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity">
                      {title}
                    </span>
                  )}
                </Link>
              ))}
            </div>
          ))}
        </div>
        <div className="stack gap-3 py-3">
          <CreditQuota />
          <UserButton size="lg" />
        </div>
      </div>
    </nav>
  );
};
