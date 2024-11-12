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
  useBillingControllerChangeSubscriptionPlan,
  useBillingControllerCreateCheckoutSession,
  useBillingControllerGetPlans,
  useOrganizationsControllerFindOne,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/use-auth";
import { ReloadIcon } from "@radix-ui/react-icons";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function BillingPageContent() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();
  const { toast } = useToast();
  const [clickedButtonIndex, setClickedButtonIndex] = useState<null | number>(
    -1
  );

  const { data: plans } = useBillingControllerGetPlans({});

  const { data: organization } = useOrganizationsControllerFindOne({
    pathParams: {
      orgname: defaultOrgname,
    },
  });

  const {
    isPending: createCheckoutSessionLoading,
    mutateAsync: createCheckoutSesseion,
  } = useBillingControllerCreateCheckoutSession({
    onError: (error) => {
      toast({
        description: error?.stack.message,
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
  } = useBillingControllerChangeSubscriptionPlan();
  const {
    isPending: cancelSubscriptionLoading,
    mutateAsync: cancelSubscription,
  } = useBillingControllerCancelSubscriptionPlan();

  return (
    <div className="flex flex-col gap-3">
      {/* New Card for Available Plans */}
      <Card>
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
                            className="flex gap-2"
                            disabled={
                              clickedButtonIndex === plans.indexOf(plan) &&
                              cancelSubscriptionLoading
                            }
                            onClick={async () => {
                              setClickedButtonIndex(plans.indexOf(plan));
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
                            {clickedButtonIndex === plans.indexOf(plan) &&
                              cancelSubscriptionLoading && (
                                <ReloadIcon className="h-5 w-5 animate-spin" />
                              )}
                            <span>Cancel Plan</span>
                          </Button>
                        ) : organization?.plan === "FREE" ? (
                          <Button
                            className="flex gap-2"
                            disabled={
                              clickedButtonIndex === plans.indexOf(plan) &&
                              createCheckoutSessionLoading
                            }
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
                            {clickedButtonIndex === plans.indexOf(plan) &&
                              createCheckoutSessionLoading && (
                                <ReloadIcon className="h-5 w-5 animate-spin" />
                              )}
                            <span>Subscribe</span>
                          </Button>
                        ) : (
                          <Button
                            className="flex gap-2"
                            disabled={
                              clickedButtonIndex === plans.indexOf(plan) &&
                              switchSubscriptionLoading
                            }
                            onClick={async () => {
                              setClickedButtonIndex(plans.indexOf(plan));
                              try {
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
                              } catch (err) {
                                toast({
                                  description: (err as any).stack.message,
                                  title: "Error",
                                  variant: "destructive",
                                });
                              }
                            }}
                            size="sm"
                          >
                            {clickedButtonIndex === plans.indexOf(plan) &&
                              switchSubscriptionLoading && (
                                <ReloadIcon className="h-5 w-5 animate-spin" />
                              )}
                            <span>Subscribe</span>
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
    </div>
  );
}
