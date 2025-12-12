import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

export const Route = createFileRoute("/_app/")({
  component: AppIndex,
});

function AppIndex(): JSX.Element {
  return (
    <div className="flex h-full flex-col items-center justify-center">
      <h1 className="font-bold text-2xl">Studio</h1>
      <p className="text-muted-foreground">Generated with Arches</p>
    </div>
  );
}
