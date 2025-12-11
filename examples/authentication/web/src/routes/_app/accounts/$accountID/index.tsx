import { Card } from "@archesai/ui";
import { createFileRoute } from "@tanstack/react-router";
import type { JSX } from "react";
import { Suspense } from "react";
import { useGetAccountSuspense } from "#lib/index";
import { getRouteMeta } from "#lib/site-utils";

export const metadata = getRouteMeta("/accounts/$accountID");

export const Route = createFileRoute("/_app/accounts/$accountID/")({
  component: AccountDetailsPage,
});

function AccountDetailsPage(): JSX.Element {
  const params = Route.useParams();
  const accountID = params.accountID;

  return (
    <div className="flex h-full w-full gap-4">
      <Card className="flex-1">
        <Suspense fallback={<div>Loading...</div>}>
          <AccountDetails accountID={accountID} />
        </Suspense>
      </Card>
    </div>
  );
}

function AccountDetails({ accountID }: { accountID: string }): JSX.Element {
  const {
    data: { data: account },
  } = useGetAccountSuspense(accountID);

  return (
    <div className="space-y-4 p-4">
      <h1 className="font-bold text-2xl">Account Details</h1>
      <dl className="grid grid-cols-2 gap-4">
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Created At
          </dt>
          <dd className="mt-1 text-sm">{String(account.createdAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Updated At
          </dt>
          <dd className="mt-1 text-sm">{String(account.updatedAt)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Access Token
          </dt>
          <dd className="mt-1 text-sm">{String(account.accessToken)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Access Token Expires At
          </dt>
          <dd className="mt-1 text-sm">
            {String(account.accessTokenExpiresAt)}
          </dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Account Identifier
          </dt>
          <dd className="mt-1 text-sm">{String(account.accountIdentifier)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            ID Token
          </dt>
          <dd className="mt-1 text-sm">{String(account.idToken)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Provider
          </dt>
          <dd className="mt-1 text-sm">{String(account.provider)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Refresh Token
          </dt>
          <dd className="mt-1 text-sm">{String(account.refreshToken)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">
            Refresh Token Expires At
          </dt>
          <dd className="mt-1 text-sm">
            {String(account.refreshTokenExpiresAt)}
          </dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">Scope</dt>
          <dd className="mt-1 text-sm">{String(account.scope)}</dd>
        </div>
        <div>
          <dt className="font-medium text-muted-foreground text-sm">User ID</dt>
          <dd className="mt-1 text-sm">{String(account.userID)}</dd>
        </div>
      </dl>
    </div>
  );
}
