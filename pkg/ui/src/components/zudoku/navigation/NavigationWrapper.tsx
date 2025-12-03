import type { PropsWithChildren } from "react";
import { useEffect, useRef } from "react";
import { cn } from "../utils";

const scrollIntoViewIfNeeded = (element: Element | null) => {
  if (!element) return;
  const rect = element.getBoundingClientRect();
  const isVisible = rect.top >= 0 && rect.bottom <= window.innerHeight;
  if (!isVisible) {
    element.scrollIntoView({ behavior: "smooth", block: "center" });
  }
};

export const NavigationWrapper = ({
  children,
  className,
}: PropsWithChildren<{
  className?: string;
}>) => {
  const navRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const active = navRef.current?.querySelector('[aria-current="page"]');
    scrollIntoViewIfNeeded(active ?? null);
  }, []);

  return (
    <div
      className="sticky top-(--header-height) grid grid-rows-[1fr_min-content] border-r lg:h-[calc(100vh-var(--header-height))]"
      data-pagefind-ignore="all"
    >
      <nav
        className={cn(
          "scrollbar hidden max-w-[calc(var(--side-nav-width)+var(--padding-nav-item))] shrink-0 flex-col overflow-y-auto ps-4 pe-3 text-sm lg:flex lg:ps-8",
          "-mx-(--padding-nav-item) scroll-pt-2 gap-1 pt-(--padding-content-top) pb-[8vh]",
          // Revert the padding/margin on the first child
          "-mt-2.5",
          className,
        )}
        ref={navRef}
        style={{
          maskImage: `linear-gradient(180deg, transparent 1%, rgba(0, 0, 0, 1) 20px, rgba(0, 0, 0, 1) 90%, transparent 99%)`,
        }}
      >
        {children}
      </nav>
    </div>
  );
};

NavigationWrapper.displayName = "NavigationWrapper";
