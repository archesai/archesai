"use client";

import { LogoSVG } from "@/components/logo-svg";
import { Button } from "@/components/ui/button";
// import { UserButton } from "@/components/user-button";
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
  const pathname = usePathname() as string;

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
    <nav className="flex flex-col text-sm font-md justify-between max-h-screen h-full opacity-100">
      <div>
        <div
          className={`flex items-center py-3 px-2.5 ${
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

          <Button
            className="hidden md:flex h-8 w-8"
            onClick={toggleSidebar}
            size="icon"
            variant="outline"
          >
            <Menu className="h-5 w-5 text-gray-500" />
          </Button>
        </div>

        {/* Render specific sidebar sections */}
        {sidebarSections.map((section) => (
          <div className="mt-4" key={section}>
            {
              <h2
                className={`text-xs px-3.5 uppercase font-semibold mb-2 whitespace-nowrap opacity-100 ${
                  isCollapsed ? "text-transparent" : "text-gray-400"
                }`}
              >
                {section.charAt(0).toUpperCase() + section.slice(1)}
              </h2>
            }

            {linksBySection[section].map(({ href, Icon, title }) => (
              <Link
                className={`flex items-center text-md font-medium gap-3 rounded-lg py-2 hover:bg-muted relative group
                ${
                  pathname === href || pathname.startsWith(href)
                    ? "bg-muted"
                    : ""
                } ${isCollapsed ? "justify-center" : "pl-[22px]"}   `}
                href={href}
                key={href}
              >
                <Icon
                  className="h-5 w-5 -translate-y-[-0.5px]"
                  strokeWidth={1.5}
                />
                <span className={`${isCollapsed ? "hidden" : "block"}`}>
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
      <div className="stack gap-3 py-3 px-3 justify-center">
        <CreditQuota />
        {/* <UserButton size="lg" /> */}
      </div>
    </nav>
  );
};
