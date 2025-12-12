import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetInvitationSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/invitations/$invitationID");

export const Route = createFileRoute("/_app/invitations/$invitationID/")({
  component: InvitationDetailsPage,
});

function InvitationDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const invitationID = params.invitationID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <InvitationDetails invitationID={invitationID} />
        </Suspense>
      </Card>
    </div>
  );
}

function InvitationDetails({
  invitationID,
}: {
  invitationID: string;
}): JSX.Element {
  const {
    data: { data: invitation },
  } = useGetInvitationSuspense(invitationID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Invitation Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(invitation.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(invitation.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Email</dt>
          <dd className="mt-1 text-sm">{String(invitation.email)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Expires At
          </dt>
          <dd className="mt-1 text-sm">{String(invitation.expiresAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Inviter ID
          </dt>
          <dd className="mt-1 text-sm">{String(invitation.inviterID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Organization ID
          </dt>
          <dd className="mt-1 text-sm">{String(invitation.organizationID)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Role</dt>
          <dd className="mt-1 text-sm">{String(invitation.role)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Status</dt>
          <dd className="mt-1 text-sm">{String(invitation.status)}</dd>
        </div>
      </dl>
    </div>
  );
}
