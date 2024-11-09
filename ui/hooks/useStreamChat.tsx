import { ContentControllerFindAllResponse } from "@/generated/archesApiComponents";
import { queryKeyFn } from "@/generated/archesApiContext";
import { ContentEntity } from "@/generated/archesApiSchemas";
import { useQueryClient } from "@tanstack/react-query";

import { useAuth } from "./useAuth";

export const useStreamChat = () => {
  const { accessToken } = useAuth();
  const queryClient = useQueryClient();

  const streamContent = (
    orgname: string,
    labelId: string,
    content: ContentEntity
  ) => {
    queryClient.setQueryData(
      queryKeyFn({
        operationId: "contentControllerFindAll",
        path: "/organizations/{orgname}/content",
        variables: {
          headers: {
            authorization: `Bearer ${accessToken}`,
          },
          pathParams: {
            orgname: orgname,
          },
          queryParams: {
            sortBy: "createdAt",
            sortDirection: "desc",
          },
        },
      }),
      (oldData: ContentControllerFindAllResponse) => {
        if (!oldData) {
          oldData = {
            metadata: { limit: 100, offset: 0, totalResults: 0 },
            results: [],
          };
        }
        const prevStreamedMessage = oldData.results?.find(
          (i) => i.id === content.id
        );
        if (prevStreamedMessage) {
          return {
            ...oldData,
            results: [
              { ...prevStreamedMessage, answer: content.text },
              ...(oldData.results || [])
                .filter((i) => i.createdAt !== prevStreamedMessage?.createdAt)
                .filter((i) => i.id !== "pending"),
            ],
          };
        } else {
          return {
            ...oldData,
            results: [content, ...(oldData.results || [])],
          };
        }
      }
    );
  };

  return {
    streamContent,
  };
};
