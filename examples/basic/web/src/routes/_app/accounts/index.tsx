import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import AccountDataTable from "#components/datatables/account-datatable";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/accounts");

export const Route = createFileRoute("/_app/accounts/")({
  component: AccountsPage,
});

function AccountsPage(): JSX.Element {
  return <AccountDataTable />;
}
