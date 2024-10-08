import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useBillingControllerGetPlans } from "@/generated/archesApiComponents";
import { PlanEntity } from "@/generated/archesApiSchemas";
import { Check } from "lucide-react";

const pricingList = [
  {
    benefitList: [
      "1 Team member",
      "2 GB Storage",
      "Upto 4 pages",
      "Community support",
      "lorem ipsum dolor",
    ],
  },
  {
    benefitList: [
      "4 Team member",
      "4 GB Storage",
      "Upto 6 pages",
      "Priority support",
      "lorem ipsum dolor",
    ],
  },
  {
    benefitList: [
      "10 Team member",
      "8 GB Storage",
      "Upto 10 pages",
      "Priority support",
      "lorem ipsum dolor",
    ],
  },
];

export const Pricing = () => {
  const { data: plans } = useBillingControllerGetPlans({});
  return (
    <section className="container py-24 sm:py-32" id="pricing">
      <h2 className="text-3xl md:text-4xl font-bold text-center">
        Get
        <span className="bg-gradient-to-b from-primary/60 to-primary text-transparent bg-clip-text">
          {" "}
          Unlimited{" "}
        </span>
        Access
      </h2>
      <h3 className="text-xl text-center text-muted-foreground pt-4 pb-8">
        Lorem ipsum dolor sit amet consectetur adipisicing elit. Alias
        reiciendis.
      </h3>
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        {plans?.toReversed().map((plan: PlanEntity, i) => (
          <Card
            className={
              plan?.metadata?.key === "STANDARD"
                ? "drop-shadow-xl shadow-black/10 dark:shadow-white/10"
                : ""
            }
            key={plan.name}
          >
            <CardHeader>
              <CardTitle className="flex item-center justify-between">
                {plan.name}
                {plan?.metadata?.key === "STANDARD" ? (
                  <Badge className="text-sm text-primary" variant="secondary">
                    Most popular
                  </Badge>
                ) : null}
              </CardTitle>
              <div>
                <span className="text-3xl font-bold">
                  ${plan.unitAmount / 100}
                </span>
                <span className="text-muted-foreground"> /month</span>
              </div>

              <CardDescription>{plan.description}</CardDescription>
            </CardHeader>

            <CardContent>
              <Button className="w-full">Choose plan</Button>
            </CardContent>

            <hr className="w-4/5 m-auto mb-4" />

            <CardFooter className="flex">
              <div className="space-y-4">
                {pricingList[i].benefitList.map((benefit: string) => (
                  <span className="flex" key={benefit}>
                    <Check className="text-green-500" />{" "}
                    <h3 className="ml-2">{benefit}</h3>
                  </span>
                ))}
              </div>
            </CardFooter>
          </Card>
        ))}
      </div>
    </section>
  );
};
