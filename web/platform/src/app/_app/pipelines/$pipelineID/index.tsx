import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/pipelines/$pipelineID/")({
  component: RouteComponent,
});

function RouteComponent() {
  return null;
}
