"use client";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader } from "@/components/ui/dialog";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  useBillingControllerCreateCheckoutSession,
  useBillingControllerGetPlans,
  useBillingControllerListPaymentMethods,
  useBillingControllerRemovePaymentMethod,
  useOrganizationsControllerFindOne,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import {
  CardElement,
  Elements,
  useElements,
  useStripe,
} from "@stripe/react-stripe-js";
import { loadStripe } from "@stripe/stripe-js";
import { useRouter } from "next/navigation";
import { useState } from "react";

const stripePromise = loadStripe("your-publishable-key"); // Replace with your Stripe publishable key

export default function BillingPageContent() {
  const router = useRouter();
  const { defaultOrgname } = useAuth();
  const [showAddCardModal, setShowAddCardModal] = useState(false);

  const { data: plans, isLoading: loadingPlans } = useBillingControllerGetPlans(
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

  const { data: organization, isLoading: organizationLoading } =
    useOrganizationsControllerFindOne({
      pathParams: {
        orgname: defaultOrgname,
      },
    });

  const { mutateAsync: deletePaymentMethod } =
    useBillingControllerRemovePaymentMethod();
  const { mutateAsync: createCheckoutSesseion } =
    useBillingControllerCreateCheckoutSession();

  const handleAddCard = () => {
    setShowAddCardModal(true);
  };

  const handleAddCardSuccess = () => {
    setShowAddCardModal(false);
  };

  return (
    <>
      <Card>
        <CardHeader>
          <CardTitle className="text-xl">Current Plan</CardTitle>
          <CardDescription>
            View your current plan and upgrade to a different plan.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {organizationLoading ? (
            <p>Loading...</p>
          ) : (
            <>
              {organization && (
                <>
                  <Input disabled value={organization.plan} />
                  <h6>
                    {plans
                      ? plans.find((p) => p.metadata?.key === organization.plan)
                          ?.description || ""
                      : ""}
                  </h6>
                </>
              )}
            </>
          )}
        </CardContent>
      </Card>

      {/* New Card for Available Plans */}
      <Card className="mt-4">
        <CardHeader>
          <CardTitle className="text-xl">Available Plans</CardTitle>
          <CardDescription>
            Subscribe to a plan to unlock additional features.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {loadingPlans ? (
            <p>Loading plans...</p>
          ) : (
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
                    {plans.map((plan) => (
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
                          <Button
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
                            variant="default"
                          >
                            Subscribe
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              ) : (
                <p>No plans available.</p>
              )}
            </>
          )}
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
                            variant="destructive"
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
              <Button onClick={handleAddCard} style={{ marginTop: "1rem" }}>
                Add New Card
              </Button>
            </>
          )}
        </CardContent>

        {showAddCardModal && (
          <Dialog
            onOpenChange={(open) => setShowAddCardModal(open)}
            open={showAddCardModal}
          >
            <DialogHeader>Add New Card</DialogHeader>
            <DialogContent>
              <Elements stripe={stripePromise}>
                <AddCardForm
                  onSuccess={handleAddCardSuccess}
                  orgname={defaultOrgname}
                />
              </Elements>
            </DialogContent>
          </Dialog>
        )}
      </Card>
    </>
  );
}

interface AddCardFormProps {
  onSuccess: (paymentMethod: any) => void;
  orgname: string;
}

const AddCardForm: React.FC<AddCardFormProps> = ({ onSuccess, orgname }) => {
  const stripe = useStripe();
  const elements = useElements();
  const [errorMessage, setErrorMessage] = useState("");
  const [processing, setProcessing] = useState(false);

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (!stripe || !elements) {
      return;
    }

    setProcessing(true);

    // Fetch the setup intent client secret from the server
    const res = await fetch(
      `/organizations/${orgname}/payment-methods/setup-intent`,
      {
        credentials: "include",
        method: "POST",
      }
    );
    const { clientSecret } = await res.json();

    const cardElement = elements.getElement(CardElement);
    if (!cardElement) {
      setErrorMessage("Card element not found");
      setProcessing(false);
      return;
    }

    const result = await stripe.confirmCardSetup(clientSecret, {
      payment_method: {
        card: cardElement,
      },
    });

    if (result.error) {
      setErrorMessage(result.error.message || "An error occurred");
      setProcessing(false);
    } else {
      // Optionally set as default payment method
      await fetch(`/organizations/${orgname}/payment-methods/default`, {
        body: JSON.stringify({
          paymentMethodId: result.setupIntent.payment_method,
        }),
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        method: "POST",
      });

      // Fetch the new payment method details
      const paymentMethodRes = await fetch(
        `/organizations/${orgname}/payment-methods/${result.setupIntent.payment_method}`,
        { credentials: "include" }
      );
      const newPaymentMethod = await paymentMethodRes.json();

      onSuccess(newPaymentMethod);
      setProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <FormField
        name="card"
        render={() => {
          return (
            <FormItem>
              <FormLabel>Card Details</FormLabel>
              <FormControl>
                <CardElement options={{ hidePostalCode: true }} />
              </FormControl>
              {errorMessage && <FormMessage>{errorMessage}</FormMessage>}
            </FormItem>
          );
        }}
      ></FormField>
      <Button disabled={!stripe || processing} type="submit">
        {processing ? "Processing..." : "Add Card"}
      </Button>
    </form>
  );
};
