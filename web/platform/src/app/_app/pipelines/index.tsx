import type { JSX } from "react"

import { createFileRoute } from "@tanstack/react-router"

import PipelineDataTable from "#components/datatables/pipeline-datatable"
import { getRouteMeta } from "#lib/site-utils"

export const metadata = getRouteMeta("/pipelines")

export const Route = createFileRoute("/_app/pipelines/")({
  component: PipelinesPage
})

export default function PipelinesPage(): JSX.Element {
  return <PipelineDataTable />
}
