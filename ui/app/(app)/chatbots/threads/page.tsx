"use client";
import { DataTable } from "@/components/datatable/data-table";
import { DataTableColumnHeader } from "@/components/datatable/data-table-column-header";
import {
  ThreadsControllerFindAllPathParams,
  ThreadsControllerRemoveVariables,
  useThreadsControllerFindAll,
  useThreadsControllerRemove,
} from "@/generated/archesApiComponents";
import { ThreadEntity } from "@/generated/archesApiSchemas";
import { useAuth } from "@/hooks/useAuth";
import { ListMinus } from "lucide-react";
import { useRouter } from "next/navigation";

export default function ChatbotsPageContent() {
  const { defaultOrgname } = useAuth();
  const router = useRouter();

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
              <div className="flex gap-2">
                <span
                  className="max-w-[500px] truncate font-medium"
                  onClick={() =>
                    router.push(`/chatbots/chat?threadId=${row.original.id}`)
                  }
                >
                  {row.original.name}
                </span>
              </div>
            );
          },
          header: ({ column }) => (
            <DataTableColumnHeader column={column} title="Name" />
          ),
        },
      ]}
      dataIcon={<ListMinus />}
      defaultView="table"
      findAllPathParams={{
        orgname: defaultOrgname,
      }}
      getDeleteVariablesFromItem={(thread) => ({
        pathParams: {
          orgname: defaultOrgname,
          threadId: thread.id,
        },
      })}
      handleSelect={(chatbot) =>
        router.push(`/chatbots/chat?threadId=${chatbot.id}`)
      }
      itemType="thread"
      useFindAll={useThreadsControllerFindAll}
      useRemove={useThreadsControllerRemove}
    />
  );
}
