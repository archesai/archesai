import { PanelLeftIcon } from "lucide-react";
import type { PropsWithChildren } from "react";
import { useState } from "react";
import { Drawer, DrawerTrigger } from "#components/shadcn/drawer";
import { useCurrentNavigation, useZudoku } from "./context/ZudokuContext";
import { Navigation } from "./navigation/Navigation";
import { Slot } from "./Slot";
import { cn } from "./utils";

const useNavigation = () => {
  return { state: "idle" as "idle" | "loading" };
};

export const Main = ({ children }: PropsWithChildren) => {
  const [isDrawerOpen, setDrawerOpen] = useState(false);
  const { navigation } = useCurrentNavigation();
  const hasNavigation = navigation.length > 0;
  const isNavigating = useNavigation().state === "loading";
  const { options } = useZudoku();

  return (
    <Drawer
      direction={options?.site?.dir === "rtl" ? "right" : "left"}
      onOpenChange={(open) => setDrawerOpen(open)}
      open={isDrawerOpen}
    >
      {hasNavigation && (
        <Navigation
          navigation={navigation}
          onRequestClose={() => setDrawerOpen(false)}
        />
      )}
      {hasNavigation && (
        <div className="sticky start-0 end-0 top-0 z-10 -mx-4 border-b bg-background/80 px-4 py-2 backdrop-blur-xs lg:hidden">
          <DrawerTrigger className="flex items-center gap-2 px-4">
            <PanelLeftIcon
              size={16}
              strokeWidth={1.5}
            />
            <span className="text-sm">Menu</span>
          </DrawerTrigger>
        </div>
      )}
      <main
        className={cn(
          "px-4 lg:px-8 lg:pe-8",
          !hasNavigation && "col-span-full",
          isNavigating && "animate-pulse",
        )}
        data-pagefind-body
      >
        <Slot.Target name="content-before" />
        {children}
        <Slot.Target name="content-after" />
      </main>
    </Drawer>
  );
};
