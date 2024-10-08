"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  BarChart2,
  ClapperboardIcon,
  Image,
  MessageCircle,
  Server,
  Upload,
} from "lucide-react";
import { useRouter } from "next/navigation";

const cardData = [
  {
    buttonText: "Go to Import",
    description: "Upload files or input URLs to import data.",
    icon: Upload,
    link: "/import/file",
    title: "Import Data",
  },
  {
    buttonText: "Analyze Data",
    description: "Analyze your data with our powerful tools.",
    icon: BarChart2,
    link: "/analyze",
    title: "Analyze Data",
  },
  {
    buttonText: "View Content",
    description: "View and manage your content.",
    icon: Server,
    link: "/content",
    title: "View Content",
  },
  {
    buttonText: "Create Chatbot",
    description: "Set up a chatbot using imported data.",
    icon: MessageCircle,
    link: "/chatbots",
    title: "Create Chatbot",
  },
  {
    buttonText: "Create Image",
    description: "Create images using generative AI.",
    icon: Image,
    link: "/images",
    title: "Create Image",
  },
  {
    buttonText: "Create Animation",
    description: "Create animations using generative AI.",
    icon: ClapperboardIcon,
    link: "/images",
    title: "Create Animation",
  },
];

export default function HomePage() {
  const router = useRouter();

  return (
    <div className="flex flex-col items-center justify-center h-full p-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {cardData.map((card, index) => (
          <Card
            className="text-center hover:shadow-lg transition-shadow"
            key={index}
          >
            <CardHeader>
              <card.icon className="mx-auto mb-2 w-12 h-12 text-gray-700 dark:text-foreground" />
              <CardTitle className="text-xl font-semibold">
                {card.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-gray-600">{card.description}</p>
            </CardContent>
            <CardFooter className="justify-center">
              <Button className="h-8" onClick={() => router.push(card.link)}>
                {card.buttonText}
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
}
