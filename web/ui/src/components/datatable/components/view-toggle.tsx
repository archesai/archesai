import type { JSX } from "react";

import { GridIcon, ListIcon } from "#components/custom/icons";
import { Button } from "#components/shadcn/button";
import { useToggleView } from "#hooks/use-toggle-view";
import { cn } from "#lib/utils";

export function ViewToggle(): JSX.Element {
  const { setView, view } = useToggleView();
  return (
    <div className="flex gap-2">
      <Button
        className={cn(
          view === "table" ? "text-primary hover:text-primary" : "",
        )}
        onClick={() => {
          setView("table");
        }}
        size={"sm"}
        variant={"ghost"}
      >
        <ListIcon />
      </Button>
      <Button
        className={cn(view === "grid" ? "text-primary hover:text-primary" : "")}
        onClick={() => {
          setView("grid");
        }}
        size={"sm"}
        variant={"ghost"}
      >
        <GridIcon />
      </Button>
    </div>
  );
}
