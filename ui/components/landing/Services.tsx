import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import Image from "next/image";

import { ChartIcon, MagnifierIcon, WalletIcon } from "./Icons";

interface ServiceProps {
  description: string;
  icon: JSX.Element;
  title: string;
}

const serviceList: ServiceProps[] = [
  {
    description:
      "Lorem ipsum dolor sit amet consectetur adipisicing elit. Nisi nesciunt est nostrum omnis ab sapiente.",
    icon: <ChartIcon />,
    title: "Code Collaboration",
  },
  {
    description:
      "Lorem ipsum dolor sit amet consectetur adipisicing elit. Nisi nesciunt est nostrum omnis ab sapiente.",
    icon: <WalletIcon />,
    title: "Project Management",
  },
  {
    description:
      "Lorem ipsum dolor sit amet consectetur adipisicing elit. Nisi nesciunt est nostrum omnis ab sapiente.",
    icon: <MagnifierIcon />,
    title: "Task Automation",
  },
];

export const Services = () => {
  return (
    <section className="container py-24 sm:py-32">
      <div className="grid place-items-center gap-8 lg:grid-cols-[1fr,1fr]">
        <div>
          <h2 className="text-3xl font-bold md:text-4xl">
            <span className="bg-gradient-to-b from-primary/60 to-primary bg-clip-text text-transparent">
              Client-Centric{" "}
            </span>
            Services
          </h2>

          <p className="mb-8 mt-4 text-xl text-muted-foreground">
            Lorem ipsum dolor sit amet consectetur, adipisicing elit. Veritatis
            dolor.
          </p>

          <div className="flex flex-col gap-8">
            {serviceList.map(({ description, icon, title }: ServiceProps) => (
              <Card key={title}>
                <CardHeader className="flex items-start justify-start gap-4 space-y-1 md:flex-row">
                  <div className="mt-1 rounded-2xl bg-primary/20 p-1">
                    {icon}
                  </div>
                  <div>
                    <CardTitle>{title}</CardTitle>
                    <CardDescription className="text-md mt-2">
                      {description}
                    </CardDescription>
                  </div>
                </CardHeader>
              </Card>
            ))}
          </div>
        </div>

        <Image
          alt="About services"
          className="w-[300px] object-contain md:w-[500px] lg:w-[600px]"
          height={300}
          src={"/cube-leg.png"}
          width={500}
        />
      </div>
    </section>
  );
};
