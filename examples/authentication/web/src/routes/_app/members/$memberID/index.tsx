import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetMemberSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/members/$memberID");

export const Route = createFileRoute("/_app/members/$memberID/")({
  component: MemberDetailsPage,
});

function MemberDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const memberID = params.memberID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <MemberDetails memberID={memberID} />
        </Suspense>
      </Card>
    </div>
  );
}

function MemberDetails({ memberID }: { memberID: string }): JSX.Element {
  const {
    data: { data: member },
  } = useGetMemberSuspense(memberID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Member Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(member.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(member.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(member.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Role</dt>
          <dd className="mt-1 text-sm">{String(member.role)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">User ID</dt>
          <dd className="mt-1 text-sm">{String(member.userID)}</dd>
        </div>
      </dl>
    </div>
  );
}
