import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import InvitationDataTable from "#components/datatables/invitation-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/invitations");

export const Route = createFileRoute("/_app/invitations/")({
  component: InvitationsPage,
});

function InvitationsPage(): JSX.Element {
  return <InvitationDataTable />;
}
