import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import ToolDataTable from "#components/datatables/tool-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/tools");

export const Route = createFileRoute("/_app/tools/")({
  component: ToolsPage,
});

export default function ToolsPage(): JSX.Element {
  return <ToolDataTable />;
}
