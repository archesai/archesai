"use client";

import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import {
  useContentControllerFindAll,
  useLabelsControllerCreate,
  usePipelinesControllerCreatePipelineRun,
} from "@/generated/archesApiComponents";
import { useAuth } from "@/hooks/useAuth";
import { useStreamChat } from "@/hooks/useStreamChat";
import { cn } from "@/lib/utils";
import { RefreshCcw } from "lucide-react";
import { useSearchParams } from "next/navigation";
import { ChangeEvent, KeyboardEvent, useEffect, useRef, useState } from "react";

export default function ChatbotChatPage() {
  const { toast } = useToast();

  const searchParams = useSearchParams();
  const [labelId, setLabelId] = useState<string>("");
  const { defaultOrgname } = useAuth();
  const [message, setMessage] = useState<string>("");
  const tid = searchParams?.get("labelId");
  const [enableFetching, setEnableFetching] = useState(false);

  useEffect(() => {
    if (tid) {
      setLabelId(tid as string);
      setEnableFetching(true);
    }
  }, [tid]);

  const { streamContent } = useStreamChat();

  const { data: messages } = useContentControllerFindAll(
    {
      pathParams: {
        orgname: defaultOrgname,
      },
      queryParams: {
        filters: labelId
          ? (JSON.stringify({
              labelId: labelId,
            }) as any)
          : undefined,
        sortBy: "createdAt",
        sortDirection: "desc",
      },
    },
    {
      enabled: enableFetching,
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

  const { mutateAsync: createLabel } = useLabelsControllerCreate();
  const { mutateAsync: createPipelineRun } =
    usePipelinesControllerCreatePipelineRun();
  const handleSend = async () => {
    if (!message.trim()) return; // Prevent sending empty messages

    let currentLabelId = labelId;
    if (!labelId) {
      try {
        const label = await createLabel({
          pathParams: {
            orgname: defaultOrgname,
          },
        });
        setLabelId(label.id);
        currentLabelId = label.id;
      } catch (error) {
        console.error("Failed to create label:", error);
        // Optionally, show a toast notification
        return;
      }
    }

    try {
      setMessage("");
      streamContent(defaultOrgname, currentLabelId, {
        createdAt: new Date().toISOString(),
        credits: 0,
        description: "Pending",
        id: "pending",
        name: "Pending",
        orgname: defaultOrgname,
        text: message.trim(),
      });
      await createPipelineRun({
        body: {
          text: message.trim(),
        },
        pathParams: {
          orgname: defaultOrgname,
          pipelineId: "pipeline-id",
        },
      });
    } catch (error) {
      toast({
        description: (error as any).stack.message,
        title: "Failed to send message",
      });
    }
  };

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const handleChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    if (e.target.value.endsWith("@")) {
      setOpen(true);
    } else {
      setOpen(false);
    }
    setMessage(e.target.value);
  };
  const [, setOpen] = useState(false);

  return (
    <div className="relative flex h-full gap-6">
      {/* Full Screen Button */}
      <div className="absolute left-0 top-0 z-10 hidden flex-col gap-2 bg-transparent md:flex">
        <Button
          className="text-muted-foreground hover:text-primary"
          onClick={() => {
            setLabelId("");
          }}
          size="icon"
          variant={"ghost"}
        >
          <RefreshCcw className="h-5 w-5" />
        </Button>
      </div>

      {/* Chat Body */}
      <div className="flex flex-1 flex-col">
        {/* Message Area */}
        <ScrollArea className="flex-1 p-4">
          <div className="flex flex-col gap-4 px-8 xl:px-52">
            {messages &&
              messages.results
                ?.slice()
                .reverse()
                .map((msg) => (
                  <div className="flex flex-col gap-2" key={msg.id}>
                    {/* User Message */}
                    <div className="flex justify-end py-2">
                      <div className="rounded-lg bg-gray-200 px-4 py-2 text-gray-800 dark:bg-gray-800 dark:text-gray-200">
                        {msg.text}
                      </div>
                    </div>
                    {/* Bot Response */}
                    <div className="flex items-start gap-2 py-2">
                      {/* <Avatar>
                        <ArchesLogo scale={0.124} size="sm" />
                      </Avatar> */}
                      {msg.id === "pending" ? (
                        <div className="flex items-center justify-center">
                          <div className="pulse h-5 w-5 rounded-full bg-black"></div>
                        </div>
                      ) : (
                        <div className="rounded-lg py-2">
                          {(msg.text as string)
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
          onSubmit={(e) => {
            e.preventDefault();
            handleSend();
          }}
        >
          <div className="flex items-center gap-2">
            <Textarea
              className="text-md max-h-40 flex-1 resize-none rounded-lg bg-background text-gray-800 focus-visible:ring-0 focus-visible:ring-transparent focus-visible:ring-offset-0 dark:text-gray-200"
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
                !message.trim() && "cursor-not-allowed opacity-50"
              )}
              disabled={!message.trim()}
              type="submit"
            >
              <svg
                className="h-5 w-5 text-white"
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
