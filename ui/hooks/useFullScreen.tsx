import { fullScreenAtom } from "@/state/fullScreenAtom";
import { useAtom } from "jotai";
import { usePathname } from "next/navigation";
import { useEffect } from "react";

import { useSidebar } from "./useSidebar";

export const useFullScreen = () => {
  const pathname = usePathname();
  const [isFullScreen, setIsFullScreen] = useAtom(fullScreenAtom);
  const { collapseSidebar, expandSidebar } = useSidebar();

  useEffect(() => {
    if (pathname !== "/chatbots/single") {
      setIsFullScreen(false);
    }
  }, [pathname]);

  const toggleFullscreen = () => {
    if (isFullScreen) {
      expandSidebar();
    } else {
      collapseSidebar();
    }
    setIsFullScreen((prev) => !prev);
  };

  const openFullScreen = () => {
    setIsFullScreen(true);
    collapseSidebar();
  };

  const closeFullScreen = () => {
    setIsFullScreen(false);
    expandSidebar();
  };

  return {
    closeFullScreen,
    isFullScreen,
    openFullScreen,
    toggleFullscreen,
  };
};
