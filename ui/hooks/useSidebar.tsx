// src/hooks/useSidebar.js

import { useAtom } from "jotai";

import { sidebarCollapsedAtom } from "../state/sidebarAtom";

/**
 * Custom hook to manage the sidebar's collapsed state.
 *
 * @returns {Object} - Contains the current state and a toggle function.
 */
export const useSidebar = () => {
  const [isCollapsed, setIsCollapsed] = useAtom(sidebarCollapsedAtom);

  /**
   * Toggles the sidebar's collapsed state.
   */
  const toggleSidebar = () => {
    setIsCollapsed((prev) => !prev);
  };

  /**
   * Sets the sidebar's state to collapsed.
   */
  const collapseSidebar = () => {
    setIsCollapsed(true);
  };

  /**
   * Sets the sidebar's state to expanded.
   */
  const expandSidebar = () => {
    setIsCollapsed(false);
  };

  return {
    collapseSidebar,
    expandSidebar,
    isCollapsed,
    toggleSidebar,
  };
};
