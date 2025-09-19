import { getSession } from "@archesai/client";
import { createServerFn } from "@tanstack/react-start";
import {
  getCookie,
  getWebRequest,
  setResponseStatus,
} from "@tanstack/react-start/server";
import { isProblem } from "./schema-validator";

const getSessionSSR = createServerFn({
  method: "GET",
}).handler(async () => {
  const { headers } = getWebRequest();
  try {
    const sessionId = getCookie("session_id");

    if (!sessionId) {
      setResponseStatus(401);
      return null;
    }

    const session = await getSession(sessionId, {
      credentials: "include",
      headers: Object.fromEntries(headers.entries()),
    });
    setResponseStatus(200);
    return session;
  } catch (error) {
    if (isProblem(error)) {
      // Log the full Problem for debugging
      console.debug("API Problem:", {
        detail: error.detail,
        instance: error.instance,
        status: error.status,
        title: error.title,
        type: error.type,
      });
      setResponseStatus(error.status || 500);
      return null;
    }

    // For unexpected errors, log the full error
    console.error("Unexpected error during session validation:", error);
    setResponseStatus(500);
    return null;
  }
});

export default getSessionSSR;
