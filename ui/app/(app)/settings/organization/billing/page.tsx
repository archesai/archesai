"use client";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useToast } from "@/components/ui/use-toast";
import {
  useBillingControllerCancelSubscriptionPlan,
  useBillingControllerCreateCheckoutSession,
  useBillingControllerGetPlans,
  useBillingControllerListPaymentMethods,
  useBillingControllerRemovePaymentMethod,
  useBillingControllerSwitchSubscriptionPlan,
  useOrganizationsControllerFindOne,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { ReloadIcon } from "@radix-ui/react-icons";
import { useRouter } from "next/navigation";

export default function BillingPageContent() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();
  const { toast } = useToast();

  const { data: plans } = useBillingControllerGetPlans(
    {},
    {
      enabled: !!defaultOrgname,
    }
  );
  const { data: paymentMethods, isLoading: loadingMethods } =
    useBillingControllerListPaymentMethods(
      {
        pathParams: { orgname: defaultOrgname },
      },
      {
        enabled: !!defaultOrgname,
      }
    );

  const { data: organization } = useOrganizationsControllerFindOne({
    pathParams: {
      orgname: defaultOrgname,
    },
  });

  const { mutateAsync: deletePaymentMethod } =
    useBillingControllerRemovePaymentMethod({
      onError: (error) => {
        toast({
          description: error?.stack.msg,
          title: "Could not delete payment method",
          variant: "destructive",
        });
      },
      onSuccess: () => {
        toast({
          description: "The payment method has been successfully deleted.",
          title: "Payment method deleted",
          variant: "default",
        });
      },
    });
  const {
    isPending: createCheckoutSessionLoading,
    mutateAsync: createCheckoutSesseion,
  } = useBillingControllerCreateCheckoutSession({
    onError: (error) => {
      toast({
        description: error?.stack.msg,
        title: "Could not create checkout session",
        variant: "destructive",
      });
    },
    onSuccess: () => {
      toast({
        description: "The checkout session has been successfully created.",
        title: "Checkout session created",
        variant: "default",
      });
    },
  });
  const {
    isPending: switchSubscriptionLoading,
    mutateAsync: switchSubscriptionPlan,
  } = useBillingControllerSwitchSubscriptionPlan();
  const {
    isPending: cancelSubscriptionLoading,
    mutateAsync: cancelSubscription,
  } = useBillingControllerCancelSubscriptionPlan();

  return (
    <>
      {/* New Card for Available Plans */}
      <Card className="mt-4">
        <CardHeader>
          <CardTitle className="text-xl">Available Plans</CardTitle>
          <CardDescription>
            Subscribe to a plan to unlock additional features.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <>
            {plans && plans.length > 0 ? (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableCell>Plan Name</TableCell>
                    <TableCell>Description</TableCell>
                    <TableCell>Price</TableCell>
                    <TableCell>Interval</TableCell>
                    <TableCell>Actions</TableCell>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {plans.toReversed().map((plan) => (
                    <TableRow key={plan.id}>
                      <TableCell>{plan.name}</TableCell>
                      <TableCell>{String(plan.description) || "-"}</TableCell>
                      <TableCell>
                        {plan.unitAmount
                          ? `$${(plan.unitAmount / 100).toFixed(2)} ${plan.currency.toUpperCase()}`
                          : "Free"}
                      </TableCell>
                      <TableCell>
                        {plan.recurring
                          ? `${plan.recurring.interval_count} ${plan.recurring.interval}(s)`
                          : "One-time"}
                      </TableCell>
                      <TableCell>
                        {organization?.plan === plan.metadata?.key ? (
                          <Button
                            disabled={cancelSubscriptionLoading}
                            onClick={async () => {
                              await cancelSubscription({
                                pathParams: {
                                  orgname: defaultOrgname,
                                },
                              });
                              toast({
                                description: "Plan canceled successfully.",
                                title: "Success",
                                variant: "default",
                              });
                            }}
                            size="sm"
                            variant="destructive"
                          >
                            {cancelSubscriptionLoading && (
                              <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                            )}
                            Cancel Plan
                          </Button>
                        ) : organization?.plan === "FREE" ? (
                          <Button
                            disabled={createCheckoutSessionLoading}
                            onClick={async () => {
                              const data = await createCheckoutSesseion({
                                pathParams: {
                                  orgname: defaultOrgname,
                                },
                                queryParams: {
                                  planId: plan.id,
                                },
                              });
                              router.push(data.url);
                            }}
                            size="sm"
                          >
                            {createCheckoutSessionLoading && (
                              <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                            )}
                            Subscribe
                          </Button>
                        ) : (
                          <Button
                            disabled={switchSubscriptionLoading}
                            onClick={async () => {
                              await switchSubscriptionPlan({
                                pathParams: {
                                  orgname: defaultOrgname,
                                },
                                queryParams: {
                                  planId: plan.id,
                                },
                              });
                              toast({
                                description: "Plan switched successfully.",
                                title: "Success",
                                variant: "default",
                              });
                            }}
                            size="sm"
                          >
                            {switchSubscriptionLoading && (
                              <ReloadIcon className="mr-2 h-4 w-4 animate-spin" />
                            )}
                            Switch Plan
                          </Button>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            ) : (
              <p>No plans available.</p>
            )}
          </>
        </CardContent>
      </Card>

      <Card className="mt-4">
        <CardHeader>
          <CardTitle className="text-xl">Payment Methods</CardTitle>
          <CardDescription>
            Manage your payment methods and subscribe to available plans.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loadingMethods ? (
            <p>Loading...</p>
          ) : (
            <>
              {paymentMethods && paymentMethods.length > 0 ? (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableCell>Brand</TableCell>
                      <TableCell>Last 4</TableCell>
                      <TableCell>Expires</TableCell>
                      <TableCell>Actions</TableCell>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {paymentMethods.map((pm) => (
                      <TableRow key={pm.id}>
                        <TableCell>{pm.card.brand}</TableCell>
                        <TableCell>{pm.card.last4}</TableCell>
                        <TableCell>
                          {pm.card.exp_month}/{pm.card.exp_year}
                        </TableCell>
                        <TableCell>
                          <Button
                            onClick={() =>
                              deletePaymentMethod({
                                pathParams: {
                                  orgname: defaultOrgname,
                                  paymentMethodId: pm.id,
                                },
                              })
                            }
                            size="sm"
                            variant="secondary"
                          >
                            Delete
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ) : (
                <p>No payment methods available.</p>
              )}
              {/* <Button
                onClick={handleAddCard}
                style={{ marginTop: "1rem" }}
                variant="default"
              >
                Add New Card
              </Button> */}
            </>
          )}
        </CardContent>
      </Card>
    </>
  );
}
