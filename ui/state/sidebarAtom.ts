import { atomWithStorage } from "jotai/utils"; // For persistence (optional)

// Atom to manage the sidebar's collapsed state with persistence
export const sidebarCollapsedAtom = atomWithStorage(
  "arches-sidebar-collapsed",
  true
);
