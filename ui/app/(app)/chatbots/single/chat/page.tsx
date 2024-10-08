// ChatbotChatPage.tsx
"use client";

import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Textarea } from "@/components/ui/textarea";
import {
  useMessagesControllerCreate,
  useMessagesControllerFindAll,
  useThreadsControllerCreate,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { useFullScreen } from "@/hooks/useFullScreen";
import { useStreamChat } from "@/hooks/useStreamChat";
import { cn } from "@/lib/utils";
import { Maximize, Minimize } from "lucide-react";
import { useSearchParams } from "next/navigation";
import { ChangeEvent, KeyboardEvent, useEffect, useRef, useState } from "react";

export default function ChatbotChatPage() {
  const [threadId, setThreadId] = useState<string>("");
  const { defaultOrgname } = useAuth();
  const [message, setMessage] = useState<string>("");
  const searchParams = useSearchParams();
  const chatbotId = searchParams?.get("chatbotId");

  const { streamChatMessage } = useStreamChat();

  const { data: messages } = useMessagesControllerFindAll(
    {
      pathParams: {
        chatbotId: chatbotId as string,
        orgname: defaultOrgname,
        threadId: threadId as string,
      },
      queryParams: {
        sortBy: "createdAt",
        sortDirection: "desc",
      },
    },
    {
      enabled: false,
    }
  );

  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({
        behavior: "smooth",
        block: "end",
      });
    }
  }, [messages]);

  const { mutateAsync: createThread } = useThreadsControllerCreate();
  const { mutateAsync: createMessage } = useMessagesControllerCreate();

  const handleSend = async () => {
    if (!message.trim()) return; // Prevent sending empty messages

    let currentThreadId = threadId;
    if (!threadId) {
      try {
        const thread = await createThread({
          pathParams: {
            chatbotId: chatbotId as string,
            orgname: defaultOrgname,
          },
        });
        setThreadId(thread.id);
        currentThreadId = thread.id;
      } catch (error) {
        console.error("Failed to create thread:", error);
        // Optionally, show a toast notification
        return;
      }
    }

    try {
      setMessage("");
      streamChatMessage(defaultOrgname, chatbotId as string, currentThreadId, {
        answer: "",
        citations: [],
        createdAt: new Date().toISOString(),
        credits: 0,
        id: "pending",
        question: message.trim(),
        threadId: currentThreadId,
      });
      await createMessage({
        body: {
          question: message.trim(),
        },
        pathParams: {
          chatbotId: chatbotId as string,
          orgname: defaultOrgname,
          threadId: currentThreadId,
        },
      });
    } catch (error) {
      console.error("Failed to send message:", error);
      // Optionally, show a toast notification
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
  };

  const { isFullScreen, toggleFullscreen } = useFullScreen();

  return (
    <div className="flex flex-col h-full relative group">
      {/* Full Screen Button */}
      <Button
        className="bg-transparent text-bg absolute top-0 left-0 z-10"
        onClick={toggleFullscreen}
        variant={"outline"}
      >
        {/* You can use an icon instead of text if preferred */}
        {isFullScreen ? <Minimize /> : <Maximize />}
      </Button>

      {/* Chat Body */}
      <div className="flex flex-col flex-1 overflow-hidden">
        {/* Message Area */}
        <ScrollArea className="flex-1 p-4">
          <div className="space-y-4 xl:pr-56 xl:pl-52 pl-8 pr-10">
            {messages &&
              messages.results
                ?.slice()
                .reverse()
                .map((msg) => (
                  <div className="flex flex-col space-y-2" key={msg.id}>
                    {/* User Message */}
                    <div className="flex justify-end py-2">
                      <div className=" bg-gray-200 dark:bg-gray-800 text-gray-800 dark:text-gray-200 rounded-lg px-4 py-2">
                        {msg.question}
                      </div>
                    </div>
                    {/* Bot Response */}
                    <div className="flex items-start space-x-2 py-2">
                      {/* <Avatar>
                        <LogoSVG scale={0.124} size="sm" />
                      </Avatar> */}
                      {msg.id === "pending" ? (
                        <div className="flex items-center justify-center">
                          <div className="pulse bg-black h-4 w-4 rounded-full"></div>
                        </div>
                      ) : (
                        <div className="rounded-lg py-2">
                          {(msg.answer as string)
                            .replaceAll(" -", "\n-")
                            .split(/(```[\s\S]*?```)/g)
                            .map((segment, index) => {
                              const replaced = segment
                                .split(/(\*\*[^*]+\*\*|`[^`]+`|\n)/g)
                                .map((part, partIndex) => {
                                  if (
                                    part.startsWith("**") &&
                                    part.endsWith("**")
                                  ) {
                                    return (
                                      <b key={partIndex}>{part.slice(2, -2)}</b>
                                    );
                                  } else if (
                                    part.startsWith("`") &&
                                    part.endsWith("`")
                                  ) {
                                    return (
                                      <span
                                        className="markdown-code"
                                        key={partIndex}
                                      >
                                        {part.slice(1, -1)}
                                      </span>
                                    );
                                  } else if (part === "\n") {
                                    return <br key={partIndex} />;
                                  } else {
                                    return part;
                                  }
                                });

                              return <span key={index}>{replaced}</span>;
                            })}
                        </div>
                      )}
                    </div>
                  </div>
                ))}
            <div ref={messagesEndRef} />
          </div>
        </ScrollArea>

        {/* Input Form */}
        <form
          className="flex-shrink-0 h-16"
          onSubmit={(e) => {
            e.preventDefault();
            handleSend();
          }}
        >
          <div className="flex items-center space-x-2">
            <Textarea
              className="flex-1 resize-none bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-lg overflow-hidden max-h-40 focus-visible:ring-0 focus-visible:ring-transparent focus-visible:ring-offset-0"
              onChange={handleChange}
              onInput={(e) => {
                const target = e.target as HTMLTextAreaElement;
                target.style.height = "auto";
                target.style.height = `${target.scrollHeight}px`;
              }}
              onKeyDown={handleKeyDown}
              placeholder="Type your message..."
              rows={1}
              // Auto-resize functionality
              style={{
                height: "auto",
                overflow: "hidden",
              }}
              value={message}
            />
            <Button
              className={cn(
                "flex items-center justify-center p-2",
                !message.trim() && "opacity-50 cursor-not-allowed"
              )}
              disabled={!message.trim()}
              type="submit"
            >
              <svg
                className="h-5 w-5 text-white"
                fill="none"
                viewBox="0 0 20 20"
                xmlns="http://www.w3.org/2000/svg"
              >
                <title>Send Message</title>
                <path
                  d="M15.44 1.68c.69-.05 1.47.08 2.13.74.66.67.8 1.45.75 2.14-.03.47-.15 1-.25 1.4l-.09.35a43.7 43.7 0 01-3.83 10.67A2.52 2.52 0 019.7 17l-1.65-3.03a.83.83 0 01.14-1l3.1-3.1a.83.83 0 10-1.18-1.17l-3.1 3.1a.83.83 0 01-.99.14L2.98 10.3a2.52 2.52 0 01.04-4.45 43.7 43.7 0 0111.02-3.9c.4-.1.92-.23 1.4-.26Z"
                  fill="currentColor"
                ></path>
              </svg>
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
