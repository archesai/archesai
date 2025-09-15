import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/runs/$runID/")({
  component: RouteComponent,
});

function RouteComponent() {
  return <></>;
}
