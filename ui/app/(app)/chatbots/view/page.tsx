"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import ChatbotForm from "@/components/forms/chatbot-form";
import { Badge } from "@/components/ui/badge";
import {
  ChatbotsControllerFindAllPathParams,
  ChatbotsControllerRemoveVariables,
  useChatbotsControllerFindAll,
  useChatbotsControllerRemove,
} from "@/generated/archesApiComponents";
import { ChatbotEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { Bot } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

export default function ChatbotsPageContent() {
  const { defaultOrgname } = useAuth();
  const router = useRouter();

  return (
    <DataTable<
      ChatbotEntity,
      ChatbotsControllerFindAllPathParams,
      ChatbotsControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <Link
                className="max-w-[500px] truncate font-medium text-primary"
                href={`/chatbots/single?chatbotId=${row.original.id}`}
              >
                {row.original.name}
              </Link>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Name" />
          ),
        },
        {
          accessorKey: "llmBase",
          cell: ({ row }) => {
            return <Badge variant="secondary">{row.original.llmBase}</Badge>;
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Language Model" />
          ),
        },
        {
          accessorKey: "description",
          cell: ({ row }) => {
            return <span className="text-sm">{row.original.description}</span>;
          },
          enableSorting: false,
          header: ({ column }) => (
            <DataTableColumnHeader
              className="text-sm"
              column={column}
              title="Description"
            />
          ),
        },
      ]}
      createForm={<ChatbotForm />}
      dataIcon={<Bot />}
      defaultView="grid"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(chatbot) => ({
        pathParams: {
          chatbotId: chatbot.id,
          orgname: defaultOrgname,
        },
      })}
      getEditFormFromItem={(chatbot) => <ChatbotForm chatbotId={chatbot.id} />}
      handleSelect={(chatbot) =>
        router.push(`/chatbots/single?chatbotId=${chatbot.id}`)
      }
      itemType="chatbot"
      useFindAll={useChatbotsControllerFindAll}
      useRemove={useChatbotsControllerRemove}
    />
  );
}
