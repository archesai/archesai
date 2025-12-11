import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetUserSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/users/$userID");

export const Route = createFileRoute("/_app/users/$userID/")({
  component: UserDetailsPage,
});

function UserDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const userID = params.userID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <UserDetails userID={userID} />
        </Suspense>
      </Card>
    </div>
  );
}

function UserDetails({ userID }: { userID: string }): JSX.Element {
  const {
    data: { data: user },
  } = useGetUserSuspense(userID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">User Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(user.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(user.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Email</dt>
          <dd className="mt-1 text-sm">{String(user.email)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Email Verified
          </dt>
          <dd className="mt-1 text-sm">{String(user.emailVerified)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Image</dt>
          <dd className="mt-1 text-sm">{String(user.image)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Name</dt>
          <dd className="mt-1 text-sm">{String(user.name)}</dd>
        </div>
      </dl>
    </div>
  );
}
