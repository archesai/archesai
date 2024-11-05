"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Image, MessageSquareIcon, Server } from "lucide-react";
import { useRouter } from "next/navigation";

const cardData = [
  {
    buttonText: "Manage Content",
    colorClass: "text-primary", // Blue for file handling and imports
    description: "View and manage your content.",
    icon: Server,
    link: "/content",
    title: "Manage Content",
  },
  {
    buttonText: "Create Image",
    colorClass: "text-primary", // Pink for creativity and AI-generated content
    description: "Create images using generative AI.",
    icon: Image,
    link: "/images",
    title: "Create Image",
  },
  {
    buttonText: "Create Chatbot",
    colorClass: "text-primary", // Purple for AI and communication-based features
    description: "Set up a chatbot using imported data.",
    icon: MessageSquareIcon,
    link: "/chatbots",
    title: "Create Chatbot",
  },
];

export default function HomePage() {
  const router = useRouter();

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
      {cardData.map((card, index) => (
        <Card
          className={`flex flex-col text-center hover:shadow-lg ${card.colorClass} justify-between p-6 transition-shadow`}
          key={index}
        >
          <CardHeader className="flex flex-col gap-2 p-0">
            <card.icon className="mx-auto h-12 w-12 opacity-80" />
            <CardTitle className="text-xl font-medium text-foreground">
              {card.title}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">{card.description}</p>
          </CardContent>
          <CardFooter className="justify-center p-0">
            <Button
              className="h-8"
              onClick={() => router.push(card.link)}
              variant={"secondary"}
            >
              {card.buttonText}
            </Button>
          </CardFooter>
        </Card>
      ))}
    </div>
  );
}
