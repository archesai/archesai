import { fullScreenAtom } from "@/state/fullScreenAtom";
import { useAtom } from "jotai";

import { useSidebar } from "./useSidebar";

export const useFullScreen = () => {
  const [isFullScreen, setIsFullScreen] = useAtom(fullScreenAtom);
  const { collapseSidebar, expandSidebar } = useSidebar();

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
