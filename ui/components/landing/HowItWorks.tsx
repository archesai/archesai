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
      "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum",
    icon: <MedalIcon />,
    title: "Accessibility",
  },
  {
    description:
      "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum",
    icon: <MapIcon />,
    title: "Community",
  },
  {
    description:
      "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum",
    icon: <PlaneIcon />,
    title: "Scalability",
  },
  {
    description:
      "Lorem ipsum dolor sit amet, consectetur adipisicing elit. Illum quas provident cum",
    icon: <GiftIcon />,
    title: "Gamification",
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
        Lorem ipsum dolor sit amet consectetur, adipisicing elit. Veritatis
        dolor pariatur sit!
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
