"use client";
import ChatbotForm from "@/components/forms/chatbot-form";
import { useSearchParams } from "next/navigation";

export default function ChatbotConfigurationPage() {
  const searchParams = useSearchParams();
  const chatbotId = searchParams.get("chatbotId");

  return <ChatbotForm chatbotId={chatbotId as string} />;
}
