import { useEffect, useState } from "react";

type ViewType = "grid" | "table";

// Cookie utility functions
const getCookie = (name: string): null | string => {
  if (typeof document === "undefined") return null;

  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) {
    return parts.pop()?.split(";").shift() ?? null;
  }
  return null;
};

const setCookie = (name: string, value: string, days = 30) => {
  if (typeof document === "undefined") return;

  const expires = new Date();
  expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
  // biome-ignore lint/suspicious/noDocumentCookie: Cookie Store API not widely supported yet
  document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/`;
};

export const useToggleView = ({
  defaultView = "table",
}: {
  defaultView?: ViewType;
} = {}): {
  setView: (newView: ViewType) => void;
  toggleView: () => void;
  view: ViewType;
} => {
  // Always start with default view to prevent SSR hydration mismatch
  const [view, setView] = useState<ViewType>(defaultView);
  const [isInitialized, setIsInitialized] = useState(false);

  // Initialize view from cookie after hydration
  useEffect(() => {
    if (!isInitialized) {
      const savedView = getCookie("viewType") as null | ViewType;

      // On mobile, always use grid view
      if (window.innerWidth <= 768) {
        setView("grid");
        setCookie("viewType", "grid");
      } else if (savedView) {
        setView(savedView);
      }

      setIsInitialized(true);
    }
  }, [isInitialized]);

  // Handle responsive behavior
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth <= 768) {
        setView("grid");
        setCookie("viewType", "grid");
      }
    };

    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);

  const setViewWrapper = (newView: ViewType) => {
    setView(newView);
    setCookie("viewType", newView);
  };

  const toggleView = () => {
    const newView = view === "grid" ? "table" : "grid";
    setView(newView);
    setCookie("viewType", newView);
  };

  return {
    setView: setViewWrapper,
    toggleView,
    view,
  };
};
