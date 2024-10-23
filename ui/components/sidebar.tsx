"use client";

import { LogoSVG } from "@/components/logo-svg";
import { Button } from "@/components/ui/button";
import { UserButton } from "@/components/user-button";
import { siteConfig } from "@/config/site";
import { useSidebar } from "@/hooks/useSidebar";
import { Menu } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

import { CreditQuota } from "./credit-quota";

export const Sidebar = () => {
  const { isCollapsed, toggleSidebar } = useSidebar();
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
    <nav className="flex flex-col text-sm font-md justify-between max-h-screen h-full">
      <div>
        <div
          className={`flex items-center py-3 px-4 ${
            isCollapsed ? "justify-center" : "justify-between"
          }`}
        >
          {!isCollapsed && (
            <div className="flex items-center gap-1">
              <LogoSVG size={"lg"} />
            </div>
          )}

          <Button
            className="hidden md:flex h-8 w-8"
            onClick={toggleSidebar}
            size="icon"
            variant="secondary"
          >
            <Menu className="h-5 w-5" />
          </Button>
        </div>

        {/* Render specific sidebar sections */}
        {sidebarSections.map((section) => (
          <div className="mt-3" key={section}>
            {
              <h2
                className={`text-xs px-4 uppercase mb-2 whitespace-nowrap ${
                  isCollapsed ? "text-transparent" : "text-gray-400"
                }`}
              >
                {section.charAt(0).toUpperCase() + section.slice(1)}
              </h2>
            }
            {linksBySection[section].map(({ href, Icon, title }) => {
              const isSelected = pathname === href || pathname.startsWith(href);
              return (
                <Link
                  className={`flex items-center text-md font-medium gap-3 ${isSelected ? "text-foreground" : "text-muted-foreground"} py-2 hover:bg-muted hover:text-foreground relative group pl-[22px] transition-all duration-200
                ${isSelected && "bg-muted"}`}
                  href={href}
                  key={href}
                >
                  <Icon
                    className={`h-5 w-5 -translate-y-[-0.5px]`}
                    strokeWidth={1.5}
                  />
                  <span className={`${isCollapsed ? "hidden" : "block"}`}>
                    {title}
                  </span>
                  {isCollapsed && (
                    <span className="absolute left-full ml-2 whitespace-nowrap bg-gray-800 text-white text-xs rounded-md px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200 z-10">
                      {title}
                    </span>
                  )}
                </Link>
              );
            })}
          </div>
        ))}
      </div>
      <div className="stack gap-3 py-3 px-2 justify-center items-center w-full">
        <CreditQuota />
        <UserButton size={isCollapsed ? "sm" : "lg"} />
      </div>
    </nav>
  );
};
