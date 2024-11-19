import Chat from "@/components/chat";
import { getMetadata } from "@/config/site";
import { Metadata } from "next";

export const metadata: Metadata = getMetadata("/chat");

export default function ChatPage() {
  return <Chat />;
}
