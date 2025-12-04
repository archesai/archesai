import { toast } from "@archesai/ui";
import { useEffect } from "react";

export interface OnlineManager {
  isOnline: () => boolean;
  subscribe: (callback: () => void) => () => void;
}

export function useOfflineIndicator(om: OnlineManager): void {
  useEffect(() => {
    return om.subscribe(() => {
      if (om.isOnline()) {
        toast.success("online", {
          duration: 2000,
          id: "ReactQuery",
        });
      } else {
        toast.error("offline", {
          duration: Infinity,
          id: "ReactQuery",
        });
      }
    });
  }, [om]);
}
