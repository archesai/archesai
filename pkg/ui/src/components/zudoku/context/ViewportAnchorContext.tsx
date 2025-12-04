import type { ReactNode } from "react";
import { createContext, useContext, useState } from "react";

interface ViewportAnchorContextType {
  activeAnchor: string | null;
  setActiveAnchor: (anchor: string | null) => void;
}

const ViewportAnchorContext = createContext<
  ViewportAnchorContextType | undefined
>(undefined);

export const ViewportAnchorProvider = ({
  children,
}: {
  children: ReactNode;
}) => {
  const [activeAnchor, setActiveAnchor] = useState<string | null>(null);

  return (
    <ViewportAnchorContext.Provider value={{ activeAnchor, setActiveAnchor }}>
      {children}
    </ViewportAnchorContext.Provider>
  );
};

export const useViewportAnchor = () => {
  const context = useContext(ViewportAnchorContext);
  if (!context) {
    return {
      activeAnchor: null,
      setActiveAnchor: () => {},
    };
  }
  return context;
};
