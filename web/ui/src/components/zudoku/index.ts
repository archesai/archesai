// Main components

// UI components
export { AnchorLink } from "./anchor-link";
export { Autocomplete } from "./autocomplete";
export { Banner } from "./Banner";
export { CategoryHeading } from "./category-header";
export { ClientOnly } from "./client-only";
export {
  useViewportAnchor,
  ViewportAnchorProvider,
} from "./context/ViewportAnchorContext";
// Context and providers
export {
  type AuthContextType,
  AuthProvider,
  NavigationProvider,
  useAuth,
  useCurrentNavigation,
  useZudoku,
  type ZudokuConfig,
  ZudokuProvider,
} from "./context/ZudokuContext";
export { Footer } from "./Footer";
export { Header } from "./header";
export { Link, NavLink } from "./Link";
export { Layout } from "./layout";
export { MobileTopNavigation } from "./MobileTopNavigation";
export { Main } from "./main";
// Navigation components
export { Navigation } from "./navigation/Navigation";
export { NavigationBadge } from "./navigation/NavigationBadge";
export { NavigationCategory } from "./navigation/NavigationCategory";
export { NavigationItem } from "./navigation/NavigationItem";
export { NavigationWrapper } from "./navigation/NavigationWrapper";
export { Toc } from "./navigation/Toc";
export { PageProgress } from "./page-progress";
export { Pagination } from "./pagination";
export { Slot } from "./Slot";
export { Search } from "./search";
export { Spinner } from "./spinner";
export { TopNavigation } from "./TopNavigation";
export { ThemeSwitch } from "./theme-switch";
export { ProseClasses, Typography } from "./typography";

// Utilities
export { cn, joinUrl, normalizeUrl, shouldShowItem } from "./utils";
