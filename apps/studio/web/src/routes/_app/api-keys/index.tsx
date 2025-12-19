import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import APIKeyDataTable from "#components/datatables/api-key-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/api-keys");

export const Route = createFileRoute("/_app/api-keys/")({
  component: APIKeysPage,
});

function APIKeysPage(): JSX.Element {
  return <APIKeyDataTable />;
}
