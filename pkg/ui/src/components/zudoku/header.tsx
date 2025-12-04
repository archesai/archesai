import { memo } from "react";
import { Button } from "#components/shadcn/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "#components/shadcn/dropdown-menu";
import { Skeleton } from "#components/shadcn/skeleton";
import { Banner } from "./Banner";
import { ClientOnly } from "./client-only";
import { useAuth, useZudoku } from "./context/ZudokuContext";
import { Link } from "./Link";
import { MobileTopNavigation } from "./MobileTopNavigation";
import { PageProgress } from "./page-progress";
import { Slot } from "./Slot";
import { Search as SearchComponent } from "./search";
import { TopNavigation } from "./TopNavigation";
import { ThemeSwitch } from "./theme-switch";
import { cn, joinUrl } from "./utils";

const RecursiveMenu = ({
  item,
}: {
  item: {
    label: string;
    path?: string;
    icon?: React.ComponentType<{
      size?: number;
      strokeWidth?: number;
      absoluteStrokeWidth?: boolean;
    }>;
    children?: (typeof item)[];
  };
}) => {
  return item.children ? (
    <DropdownMenuSub key={item.label}>
      <DropdownMenuSubTrigger>{item.label}</DropdownMenuSubTrigger>
      <DropdownMenuPortal>
        <DropdownMenuSubContent>
          {item.children.map((child) => (
            <RecursiveMenu
              item={child}
              key={child.label}
            />
          ))}
        </DropdownMenuSubContent>
      </DropdownMenuPortal>
    </DropdownMenuSub>
  ) : (
    <Link to={item.path ?? ""}>
      <DropdownMenuItem
        className="flex gap-2"
        key={item.label}
      >
        {item.icon && (
          <item.icon
            absoluteStrokeWidth
            size={16}
            strokeWidth={1}
          />
        )}
        {item.label}
      </DropdownMenuItem>
    </Link>
  );
};

export const Header = memo(function HeaderInner() {
  const auth = useAuth();
  const { isAuthenticated, profile, isAuthEnabled } = useAuth();
  const context = useZudoku();
  const { site, options } = context;

  const accountItems = [] as Array<{
    label: string;
    path?: string;
    icon?: React.ComponentType<{
      size?: number;
      strokeWidth?: number;
      absoluteStrokeWidth?: boolean;
    }>;
    children?: Array<{
      label: string;
      path?: string;
      icon?: React.ComponentType<{
        size?: number;
        strokeWidth?: number;
        absoluteStrokeWidth?: boolean;
      }>;
    }>;
    category?: string;
  }>;

  const logoLightSrc = site?.logo
    ? /https?:\/\//.test(site.logo.src.light)
      ? site.logo.src.light
      : joinUrl(options?.basePath, site.logo.src.light)
    : undefined;
  const logoDarkSrc = site?.logo
    ? /https?:\/\//.test(site.logo.src.dark)
      ? site.logo.src.dark
      : joinUrl(options?.basePath, site.logo.src.dark)
    : undefined;

  const borderBottom = "inset-shadow-[0_-1px_0_0_var(--border)]";

  return (
    <header
      className="sticky z-10 w-full bg-background/80 backdrop-blur lg:top-0"
      data-pagefind-ignore="all"
    >
      <Banner />
      <div className={cn(borderBottom, "relative")}>
        <PageProgress />
        <div className="mx-auto flex h-(--top-header-height) max-w-screen-2xl items-center justify-between border-transparent px-4 lg:px-8">
          <div className="flex">
            <Link to={site?.logo?.href ?? "/"}>
              <div className="flex items-center gap-3.5">
                {site?.logo ? (
                  <>
                    <img
                      alt={site.logo.alt ?? site.title}
                      className="max-h-(--top-header-height) dark:hidden"
                      loading="lazy"
                      src={logoLightSrc}
                      style={{ width: site.logo.width }}
                    />
                    <img
                      alt={site.logo.alt ?? site.title}
                      className="hidden max-h-(--top-header-height) dark:block"
                      loading="lazy"
                      src={logoDarkSrc}
                      style={{ width: site.logo.width }}
                    />
                  </>
                ) : (
                  <span className="font-semibold text-2xl">{site?.title}</span>
                )}
              </div>
            </Link>
          </div>

          <div className="pointer-events-none absolute inset-x-0 hidden w-full items-center justify-center lg:flex">
            <SearchComponent className="pointer-events-auto" />
          </div>

          <div className="flex items-center gap-8">
            <MobileTopNavigation />
            <div className="hidden items-center gap-2 justify-self-end text-sm lg:flex">
              <Slot.Target name="head-navigation-start" />
              {isAuthEnabled && (
                <ClientOnly
                  fallback={<Skeleton className="mr-4 h-5 w-24 rounded-sm" />}
                >
                  {!isAuthenticated ? (
                    <Button
                      onClick={() => auth.login()}
                      variant="ghost"
                    >
                      Login
                    </Button>
                  ) : (
                    accountItems.length > 0 && (
                      <DropdownMenu modal={false}>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost">
                            {profile?.name ?? "My Account"}
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent className="w-56">
                          <DropdownMenuLabel>
                            {profile?.name ? `${profile.name}` : "My Account"}
                            {profile?.email && (
                              <div className="font-normal text-muted-foreground">
                                {profile.email}
                              </div>
                            )}
                          </DropdownMenuLabel>
                          {accountItems.filter((i) => i.category === "top")
                            .length > 0 && <DropdownMenuSeparator />}
                          {accountItems
                            .filter((i) => i.category === "top")
                            .map((i) => (
                              <RecursiveMenu
                                item={i}
                                key={i.label}
                              />
                            ))}
                          {accountItems.filter(
                            (i) => !i.category || i.category === "middle",
                          ).length > 0 && <DropdownMenuSeparator />}
                          {accountItems
                            .filter(
                              (i) => !i.category || i.category === "middle",
                            )
                            .map((i) => (
                              <RecursiveMenu
                                item={i}
                                key={i.label}
                              />
                            ))}
                          {accountItems.filter((i) => i.category === "bottom")
                            .length > 0 && <DropdownMenuSeparator />}
                          {accountItems
                            .filter((i) => i.category === "bottom")
                            .map((i) => (
                              <RecursiveMenu
                                item={i}
                                key={i.label}
                              />
                            ))}
                        </DropdownMenuContent>
                      </DropdownMenu>
                    )
                  )}
                </ClientOnly>
              )}
              <Slot.Target name="head-navigation-end" />
              <ThemeSwitch />
            </div>
          </div>
        </div>
      </div>
      <div className={cn("hidden lg:block", borderBottom)}>
        <div className="relative mx-auto max-w-screen-2xl border-transparent">
          <Slot.Target name="top-navigation-before" />
          <TopNavigation />
          <Slot.Target name="top-navigation-after" />
        </div>
      </div>
    </header>
  );
});
