"use client";

import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import { Badge } from "@/components/ui/badge";
import {
  ThreadsControllerRemoveVariables,
  useThreadsControllerFindAll,
  useThreadsControllerRemove,
} from "@/generated/archesApiComponents";
import { ThreadEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { useFilterItems } from "@/hooks/useFilterItems";
import { useSelectItems } from "@/hooks/useSelectItems";
import { endOfDay } from "date-fns";
import { User } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";

function ChatbotThreadsPage() {
  const { defaultOrgname } = useAuth();
  const router = useRouter();
  const { limit, page, range } = useFilterItems();
  const searchParams = useSearchParams();
  const chatbotId = searchParams?.get("chatbotId");

  const {
    data: threads,
    isLoading,
    isPlaceholderData,
  } = useThreadsControllerFindAll({
    pathParams: {
      chatbotId: chatbotId as string,
      orgname: defaultOrgname,
    },
    queryParams: {
      endDate: endOfDay(range.to || new Date()).toISOString(),
      limit,
      offset: page * limit,
      sortBy: "createdAt",
      sortDirection: "asc" as const,
      startDate: range.from?.toISOString(),
    },
  });
  const loading = isPlaceholderData || isLoading;
  const { mutateAsync: deleteChatbot } = useThreadsControllerRemove();

  const { selectedItems } = useSelectItems({ items: threads?.results || [] });

  return (
    <DataTable<ThreadEntity, ThreadsControllerRemoveVariables>
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
      data={threads as any}
      dataIcon={<User size={24} />}
      defaultView="table"
      deleteItem={deleteChatbot}
      getDeleteVariablesFromItem={(thread) => [
        {
          pathParams: {
            chatbotId: chatbotId as string,
            orgname: defaultOrgname,
            threadId: thread.id,
          },
        },
      ]}
      handleSelect={(chatbot) =>
        router.push(`/chatbots/single/chat?chatbotId=${chatbot.id}`)
      }
      itemType="thread"
      loading={loading}
      mutationVariables={selectedItems.map((id) => ({
        pathParams: {
          chatbotId: chatbotId as string,
          orgname: defaultOrgname,
          threadId: id,
        },
      }))}
    />
  );
}

export default ChatbotThreadsPage;
