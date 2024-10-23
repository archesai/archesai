"use client";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Image, MessageSquareIcon, Server, Upload } from "lucide-react";
import { useRouter } from "next/navigation";

const cardData = [
  {
    buttonText: "Go to Import",
    colorClass: "text-blue-600", // Blue for file handling and imports
    description: "Upload files or input URLs to import data.",
    icon: Upload,
    link: "/import",
    title: "Import Data",
  },
  {
    buttonText: "View Content",
    colorClass: "text-green-600", // Green for managing and overseeing data
    description: "View and manage your content.",
    icon: Server,
    link: "/content",
    title: "View Content",
  },
  {
    buttonText: "Create Chatbot",
    colorClass: "text-purple-600", // Purple for AI and communication-based features
    description: "Set up a chatbot using imported data.",
    icon: MessageSquareIcon,
    link: "/chatbots",
    title: "Create Chatbot",
  },
  {
    buttonText: "Create Image",
    colorClass: "text-pink-600", // Pink for creativity and AI-generated content
    description: "Create images using generative AI.",
    icon: Image,
    link: "/images",
    title: "Create Image",
  },
];

export default function HomePage() {
  const router = useRouter();

  return (
    <div className="flex flex-col justify-start h-full">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {cardData.map((card, index) => (
          <Card
            className={`flex flex-col text-center hover:shadow-lg ${card.colorClass} transition-shadow justify-between`}
            key={index}
          >
            <CardHeader>
              <card.icon className="mx-auto mb-2 w-12 h-12" />
              <CardTitle className="text-xl font-semibold text-foreground">
                {card.title}
              </CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-muted-foreground">{card.description}</p>
            </CardContent>
            <CardFooter className="justify-center">
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
    </div>
  );
}
