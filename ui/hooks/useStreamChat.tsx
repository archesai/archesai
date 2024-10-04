import { MessagesControllerFindAllResponse } from "@/generated/archesApiComponents";
import { queryKeyFn } from "@/generated/archesApiContext";
import { MessageEntity } from "@/generated/archesApiSchemas";
import { useQueryClient } from "@tanstack/react-query";

import { useAuth } from "./useAuth";

export const useStreamChat = () => {
  const { accessToken } = useAuth();
  const queryClient = useQueryClient();

  const streamChatMessage = (
    orgname: string,
    chatbotId: string,
    threadId: string,
    message: MessageEntity
  ) => {
    queryClient.setQueryData(
      queryKeyFn({
        operationId: "messagesControllerFindAll",
        path: "/organizations/{orgname}/chatbots/{chatbotId}/threads/{threadId}/messages",
        variables: {
          headers: {
            authorization: `Bearer ${accessToken}`,
          },
          pathParams: {
            chatbotId: chatbotId,
            orgname: orgname,
            threadId: threadId,
          },
          queryParams: {
            sortBy: "createdAt",
            sortDirection: "desc",
          },
        },
      }),
      (oldData: MessagesControllerFindAllResponse) => {
        if (!oldData) {
          oldData = {
            metadata: { limit: 100, offset: 0, totalResults: 0 },
            results: [],
          };
        }
        const oldMessage = oldData.results?.find((i) => i.id === message.id);
        if (oldMessage) {
          return {
            ...oldData,
            results: [
              { ...oldMessage, answer: message.answer },
              ...(oldData.results || []).filter(
                (i) => i.createdAt !== oldMessage?.createdAt
              ),
            ],
          };
        } else {
          return {
            ...oldData,
            results: [
              message,
              ...(oldData.results || []).filter((i) => i.id != "pending"),
            ],
          };
        }
      }
    );
  };

  return {
    streamChatMessage,
  };
};
