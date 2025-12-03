import type { JSX } from "react";

import { GridIcon, ListIcon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import { cn } from "#lib/utils";

interface ViewToggleProps {
  view: "table" | "grid";
  onToggle: () => void;
}

export function ViewToggle({ view, onToggle }: ViewToggleProps): JSX.Element {
  return (
    <div className="flex gap-2">
      <Button
        className={cn(
          view === "table" ? "text-primary hover:text-primary" : "",
        )}
        onClick={() => {
          if (view !== "table") onToggle();
        }}
        size={"sm"}
        variant={"ghost"}
      >
        <ListIcon />
      </Button>
      <Button
        className={cn(view === "grid" ? "text-primary hover:text-primary" : "")}
        onClick={() => {
          if (view !== "grid") onToggle();
        }}
        size={"sm"}
        variant={"ghost"}
      >
        <GridIcon />
      </Button>
    </div>
  );
}
