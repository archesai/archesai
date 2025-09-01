import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/pipelines/$pipelineId/')({
  component: RouteComponent
})

function RouteComponent() {
  return <></>
}
