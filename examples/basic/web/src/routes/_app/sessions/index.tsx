import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import SessionDataTable from "#components/datatables/session-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/sessions");

export const Route = createFileRoute("/_app/sessions/")({
  component: SessionsPage,
});

function SessionsPage(): JSX.Element {
  return <SessionDataTable />;
}
