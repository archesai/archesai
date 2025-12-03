import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import RunDataTable from "#components/datatables/run-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/runs");

export const Route = createFileRoute("/_app/runs/")({
  component: RunsPage,
});

function RunsPage(): JSX.Element {
  return <RunDataTable />;
}
