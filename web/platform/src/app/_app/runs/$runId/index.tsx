import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_app/runs/$runId/")({
  component: RouteComponent
})

function RouteComponent() {
  return <></>
}
