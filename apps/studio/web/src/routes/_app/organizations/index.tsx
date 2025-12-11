import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import OrganizationDataTable from "#components/datatables/organization-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/organizations");

export const Route = createFileRoute("/_app/organizations/")({
  component: OrganizationsPage,
});

function OrganizationsPage(): JSX.Element {
  return <OrganizationDataTable />;
}
