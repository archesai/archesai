import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetOrganizationSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/organizations/$organizationID");

export const Route = createFileRoute("/_app/organizations/$organizationID/")({
  component: OrganizationDetailsPage,
});

function OrganizationDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const organizationID = params.organizationID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <OrganizationDetails organizationID={organizationID} />
        </Suspense>
      </Card>
    </div>
  );
}

function OrganizationDetails({
  organizationID,
}: {
  organizationID: string;
}): JSX.Element {
  const {
    data: { data: organization },
  } = useGetOrganizationSuspense(organizationID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Organization Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(organization.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(organization.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Billing Email
          </dt>
          <dd className="mt-1 text-sm">{String(organization.billingEmail)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Credits</dt>
          <dd className="mt-1 text-sm">{String(organization.credits)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Logo</dt>
          <dd className="mt-1 text-sm">{String(organization.logo)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(organization.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Plan</dt>
          <dd className="mt-1 text-sm">{String(organization.plan)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Slug</dt>
          <dd className="mt-1 text-sm">{String(organization.slug)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Stripe Customer Identifier
          </dt>
          <dd className="mt-1 text-sm">
            {String(organization.stripeCustomerIdentifier)}
          </dd>
        </div>
      </dl>
    </div>
  );
}
