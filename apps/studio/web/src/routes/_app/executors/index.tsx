import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import ExecutorDataTable from "#components/datatables/executor-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/executors");

export const Route = createFileRoute("/_app/executors/")({
  component: ExecutorsPage,
});

function ExecutorsPage(): JSX.Element {
  return <ExecutorDataTable />;
}
