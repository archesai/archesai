import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import UserDataTable from "#components/datatables/user-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/users");

export const Route = createFileRoute("/_app/users/")({
  component: UsersPage,
});

function UsersPage(): JSX.Element {
  return <UserDataTable />;
}
