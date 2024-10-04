"use client";
import ChatbotForm from "@/components/forms/chatbot-form";
import { useParams } from "next/navigation";

export default function ChatbotConfigurationPage() {
  const { chatbotId } = useParams();

  return <ChatbotForm chatbotId={chatbotId as string} />;
}
