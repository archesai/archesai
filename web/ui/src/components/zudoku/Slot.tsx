import type { ReactNode } from "react";
import { createContext, useContext } from "react";

interface SlotContextType {
  slots: Record<string, ReactNode>;
}

const SlotContext = createContext<SlotContextType>({ slots: {} });

export const SlotProvider = ({
  children,
  slots = {},
}: {
  children: ReactNode;
  slots?: Record<string, ReactNode>;
}) => {
  return (
    <SlotContext.Provider value={{ slots }}>{children}</SlotContext.Provider>
  );
};

const Target = ({ name }: { name: string }) => {
  const { slots } = useContext(SlotContext);
  return <>{slots[name]}</>;
};

const Fill = ({ children }: { name: string; children: ReactNode }) => {
  return <>{children}</>;
};

export const Slot = {
  Fill,
  Provider: SlotProvider,
  Target,
};
