import type { JSX } from "react";

import { createFileRoute } from "@tanstack/react-router";

import MemberDataTable from "#components/datatables/member-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/organization/members");

export const Route = createFileRoute("/_app/organization/members/")({
  component: MembersPage,
});

export default function MembersPage(): JSX.Element {
  return <MemberDataTable />;
}
