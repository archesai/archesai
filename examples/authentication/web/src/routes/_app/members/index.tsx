import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import MemberDataTable from "#components/datatables/member-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/members");

export const Route = createFileRoute("/_app/members/")({
  component: MembersPage,
});

function MembersPage(): JSX.Element {
  return <MemberDataTable />;
}
