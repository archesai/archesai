import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { GiftIcon, MapIcon, MedalIcon, PlaneIcon } from "./Icons";

interface FeatureProps {
  description: string;
  icon: JSX.Element;
  title: string;
}

const features: FeatureProps[] = [
  {
    description:
      "Easily sign up or log in to the Arches AI platform to access powerful AI tools designed to optimize your business processes.",
    icon: <MedalIcon />,
    title: "Sign Up & Onboarding",
  },
  {
    description:
      "Integrate your data seamlessly with our platform, allowing AI to process, analyze, and deliver insights tailored to your needs.",
    icon: <MapIcon />,
    title: "Data Integration",
  },
  {
    description:
      "Scale your operations effortlessly with AI-driven automation, enabling you to handle increased workloads without losing efficiency.",
    icon: <PlaneIcon />,
    title: "AI Automation & Scalability",
  },
  {
    description:
      "Monitor performance, receive insights, and continually improve through feedback loops built into the system.",
    icon: <GiftIcon />,
    title: "Optimization & Insights",
  },
];

export const HowItWorks = () => {
  return (
    <section className="container text-center py-24 sm:py-32" id="howItWorks">
      <h2 className="text-3xl md:text-4xl font-bold ">
        How It{" "}
        <span className="bg-gradient-to-b from-primary/60 to-primary text-transparent bg-clip-text">
          Works{" "}
        </span>
        Step-by-Step Guide
      </h2>
      <p className="md:w-3/4 mx-auto mt-4 mb-8 text-xl text-muted-foreground">
        Discover how Arches AI empowers your business with advanced AI
        solutions. Here’s a simple guide on how to get started.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
        {features.map(({ description, icon, title }: FeatureProps) => (
          <Card className="bg-muted/50" key={title}>
            <CardHeader>
              <CardTitle className="grid gap-4 place-items-center">
                {icon}
                {title}
              </CardTitle>
            </CardHeader>
            <CardContent>{description}</CardContent>
          </Card>
        ))}
      </div>
    </section>
  );
};
