interface TocEntry {
  id: string;
  value: string;
  depth: number;
  children?: TocEntry[];
}

import { ListTreeIcon } from "lucide-react";
import type { CSSProperties, PropsWithChildren } from "react";
import { useEffect, useRef, useState } from "react";
import { useViewportAnchor } from "../context/ViewportAnchorContext";
import { Link } from "../Link";
import { cn } from "../utils";

const DATA_ANCHOR_ATTR = "data-active";

const TocItem = ({
  item,
  children,
  className,
  isActive,
}: PropsWithChildren<{
  item: TocEntry;
  isActive: boolean;
  className?: string;
}>) => {
  return (
    <li
      className={cn("truncate", className)}
      title={item.value}
    >
      <Link
        href={`#${item.id}`}
        {...{ [DATA_ANCHOR_ATTR]: item.id }}
        className={cn(
          isActive
            ? "text-primary"
            : "text-muted-foreground hover:text-accent-foreground",
        )}
      >
        {item.value}
      </Link>
      {children}
    </li>
  );
};

export const Toc = ({ entries }: { entries: TocEntry[] }) => {
  const { activeAnchor } = useViewportAnchor();
  const listWrapperRef = useRef<HTMLUListElement>(null);
  const paintedOnce = useRef(false);
  const [indicatorStyle, setIndicatorStyles] = useState<CSSProperties>({
    opacity: 0,
    top: 0,
  });

  // synchronize active anchor indicator with the scroll position
  useEffect(() => {
    if (!listWrapperRef.current) return;

    const activeElement = listWrapperRef.current.querySelector(
      `[${DATA_ANCHOR_ATTR}='${activeAnchor}']`,
    );

    if (!activeElement) {
      setIndicatorStyles({ opacity: 0, top: 0 });
      return;
    }

    const topParent = listWrapperRef.current.getBoundingClientRect().top;
    const topElement = activeElement.getBoundingClientRect().top;

    setIndicatorStyles({
      opacity: 1,
      top: `${topElement - topParent}px`,
    });

    if (paintedOnce.current) return;

    // after all is painted, the indicator should animate
    requestIdleCallback(() => {
      paintedOnce.current = true;
    });
  }, [activeAnchor]);

  return (
    <aside className="scrollbar sticky top-8 h-[calc(100vh-var(--header-height))] overflow-y-auto ps-1 pt-(--padding-content-top) pb-(--padding-content-bottom) text-sm lg:top-(--header-height)">
      <div className="mb-2 flex items-center gap-2 font-medium">
        <ListTreeIcon size={16} />
        On this page
      </div>
      <div className="relative ms-2 ps-4">
        <div className="absolute inset-0 end-auto w-[2px] bg-border" />
        <div
          className={cn(
            "absolute -start-px h-6 w-[4px] -translate-y-1 rounded-sm bg-primary",
            paintedOnce.current &&
              "ease-out [transition:top_150ms,opacity_325ms]",
          )}
          style={indicatorStyle}
        />
        <ul
          className="relative list-none space-y-2 font-medium"
          ref={listWrapperRef}
        >
          {entries.map((item) => (
            <TocItem
              className="ps-0"
              isActive={item.id === activeAnchor}
              item={item}
              key={item.id}
            >
              {item.children && (
                <ul className="list-none space-y-2 ps-4 pt-2">
                  {item.children.map((child) => (
                    <TocItem
                      isActive={child.id === activeAnchor}
                      item={child}
                      key={child.id}
                    />
                  ))}
                </ul>
              )}
            </TocItem>
          ))}
        </ul>
      </div>
    </aside>
  );
};
