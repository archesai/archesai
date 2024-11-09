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
    labelId: string,
    message: MessageEntity
  ) => {
    queryClient.setQueryData(
      queryKeyFn({
        operationId: "messagesControllerFindAll",
        path: "/organizations/{orgname}/labels/{labelId}/messages",
        variables: {
          headers: {
            authorization: `Bearer ${accessToken}`,
          },
          pathParams: {
            orgname: orgname,
            labelId: labelId,
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
        const prevStreamedMessage = oldData.results?.find(
          (i) => i.id === message.id
        );
        if (prevStreamedMessage) {
          return {
            ...oldData,
            results: [
              { ...prevStreamedMessage, answer: message.answer },
              ...(oldData.results || [])
                .filter((i) => i.createdAt !== prevStreamedMessage?.createdAt)
                .filter((i) => i.id !== "pending"),
            ],
          };
        } else {
          return {
            ...oldData,
            results: [message, ...(oldData.results || [])],
          };
        }
      }
    );
  };

  return {
    streamChatMessage,
  };
};
