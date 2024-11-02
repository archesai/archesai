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
  const sidebarSections = ["home", "data", "manage", "studio", "settings"];

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
    <nav className="font-md flex h-full max-h-screen flex-col justify-between p-3 text-sm">
      {/* TOP PART */}
      <div className="flex flex-col gap-3">
        {/* Render logo and collapse button */}
        <div
          className={`flex items-center ${
            isCollapsed ? "justify-center" : "justify-between"
          }`}
        >
          {!isCollapsed && <LogoSVG size={"lg"} />}

          <Button
            className="hidden h-8 w-8 md:flex"
            onClick={toggleSidebar}
            size="icon"
            variant="secondary"
          >
            <Menu className="h-5 w-5" />
          </Button>
        </div>

        {/* Render specific sidebar sections */}
        {sidebarSections.map((section) => (
          <div className="flex flex-col gap-2" key={section}>
            <h2
              className={`whitespace-nowrap text-xs uppercase ${
                isCollapsed ? "text-transparent" : "text-gray-400"
              }`}
            >
              {section.charAt(0).toUpperCase() + section.slice(1)}
            </h2>
            {linksBySection[section].map(({ href, Icon, title }) => {
              return (
                <Link
                  className={`text-md group relative flex items-center rounded-lg p-2 font-medium duration-200 hover:bg-muted hover:text-foreground ${pathname.startsWith(href) ? "bg-secondary text-secondary-foreground" : "text-muted-foreground"} `}
                  href={href}
                  key={href}
                >
                  <Icon className="h-5 w-5 flex-shrink-0" strokeWidth={1.5} />
                  <div
                    className={`ml-3 transition-opacity duration-500 ${isCollapsed ? "w-0 opacity-0" : "w-auto opacity-100"} `}
                    style={{ overflow: "hidden" }}
                  >
                    {title}
                  </div>
                  {/* Render the title as a tooltip when collapsed */}
                  {isCollapsed && (
                    <span className="absolute left-full z-10 m-2 whitespace-nowrap rounded-md bg-primary p-2 text-xs text-primary-foreground opacity-0 transition-opacity duration-200 group-hover:opacity-100">
                      {title}
                    </span>
                  )}
                </Link>
              );
            })}
          </div>
        ))}
      </div>

      {/* BOTTOM PART */}
      <div className="flex flex-col items-center gap-3">
        <CreditQuota />
        <UserButton size={isCollapsed ? "sm" : "lg"} />
      </div>
    </nav>
  );
};
