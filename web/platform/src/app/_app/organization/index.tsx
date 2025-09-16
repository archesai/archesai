import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";

import OrganizationForm from "#components/forms/organization-form";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/organization");

export const Route = createFileRoute("/_app/organization/")({
  component: OrganizationPage,
});

function OrganizationPage(): JSX.Element {
  return <OrganizationForm />;
}
