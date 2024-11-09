import { useQueryClient } from "@tanstack/react-query";
import { useEffect } from "react";
import { io } from "socket.io-client";

import { useAuth } from "./useAuth";
import { useStreamChat } from "./useStreamChat";

export const useWebsockets = ({
  overrideToken,
}: {
  overrideToken?: string;
}) => {
  const { accessToken } = useAuth();
  const queryClient = useQueryClient();

  const { streamChatMessage } = useStreamChat();

  useEffect(() => {
    if (accessToken) {
      const websocket = io(process.env.NEXT_PUBLIC_WEBSOCKET_URL as string, {
        auth: {
          token: overrideToken || accessToken,
        },
        extraHeaders: {
          Authorization: `Bearer ${overrideToken || accessToken}`,
        },
        reconnection: true,
        reconnectionAttempts: Infinity,
        reconnectionDelay: 1000,
        reconnectionDelayMax: 5000,
        transports: ["websocket"],
        withCredentials: true,
      });

      websocket.on("connect", () => {
        console.debug("connected");
      });

      websocket.on("ping", () => {});

      websocket.on("update", async (event) => {
        await queryClient.invalidateQueries({
          queryKey: event.queryKey,
        });
      });

      websocket.on("chat", (event) => {
        streamChatMessage(
          event.orgname,
          event.chatbotId,
          event.labelId,
          event.message
        );
      });

      return () => {
        websocket.close();
      };
    }
  }, [queryClient, accessToken, overrideToken]);
};
