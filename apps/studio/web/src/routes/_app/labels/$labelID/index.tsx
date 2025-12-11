import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetLabelSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/labels/$labelID");

export const Route = createFileRoute("/_app/labels/$labelID/")({
  component: LabelDetailsPage,
});

function LabelDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const labelID = params.labelID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <LabelDetails labelID={labelID} />
        </Suspense>
      </Card>
    </div>
  );
}

function LabelDetails({ labelID }: { labelID: string }): JSX.Element {
  const {
    data: { data: label },
  } = useGetLabelSuspense(labelID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Label Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(label.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(label.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(label.name)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(label.organizationID)}</dd>
        </div>
      </dl>
    </div>
  );
}
