"use client";

import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { Badge } from "@/components/ui/badge";
import {
  ThreadsControllerFindAllPathParams,
  ThreadsControllerRemoveVariables,
  useThreadsControllerFindAll,
  useThreadsControllerRemove,
} from "@/generated/archesApiComponents";
import { ThreadEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { User } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";

function ChatbotThreadsPage() {
  const { defaultOrgname } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const chatbotId = searchParams?.get("chatbotId");

  return (
    <DataTable<
      ThreadEntity,
      ThreadsControllerFindAllPathParams,
      ThreadsControllerRemoveVariables
    >
      columns={[
        {
          accessorKey: "name",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="max-w-[500px] truncate font-medium">
                  {row.original.name}
                </span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Name" />
          ),
        },
        {
          accessorKey: "createdAt",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <span className="font-medium">{row.original.createdAt}</span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Created" />
          ),
        },
        {
          accessorKey: "numMessages",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <Badge variant="outline">{row.original.numMessages}</Badge>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Messages" />
          ),
        },
        {
          accessorKey: "credits",
          cell: ({ row }) => {
            return (
              <div className="flex space-x-2">
                <Badge variant="outline">{row.original.credits}</Badge>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Credits" />
          ),
        },
      ]}
      content={() => (
        <div className="flex w-full justify-center items-center h-full">
          <User size={"100px"} />
        </div>
      )}
      dataIcon={<User size={24} />}
      defaultView="table"
      findAllPathParams={{
        chatbotId: chatbotId as string,
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(thread) => ({
        pathParams: {
          chatbotId: chatbotId as string,
          orgname: defaultOrgname,
          threadId: thread.id,
        },
      })}
      handleSelect={(chatbot) =>
        router.push(`/chatbots/single/chat?chatbotId=${chatbot.id}`)
      }
      itemType="thread"
      useFindAll={useThreadsControllerFindAll}
      useRemove={useThreadsControllerRemove}
    />
  );
}

export default ChatbotThreadsPage;
