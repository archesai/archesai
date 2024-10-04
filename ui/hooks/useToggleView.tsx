import { viewAtom } from "@/state/viewAtom";
import { useAtom } from "jotai"; // Adjust path as necessary
import { useEffect, useState } from "react";

export const useToggleView = () => {
  const [view, setView] = useAtom(viewAtom);
  const [width, setWidth] = useState(window.innerWidth);

  const toggleView = () => {
    setView((prev) => (prev === "grid" ? "table" : "grid"));
  };

  useEffect(() => {
    const handleResize = () => {
      setWidth(window.innerWidth);
    };

    window.addEventListener("resize", handleResize);

    if (width <= 768) {
      setView("grid");
    }

    return () => window.removeEventListener("resize", handleResize);
  }, [width, setView]);

  return {
    setView,
    toggleView,
    view,
  };
};
