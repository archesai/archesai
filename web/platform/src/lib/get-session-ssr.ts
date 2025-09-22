import { getSession, isProblem } from "@archesai/client";
import { createServerFn } from "@tanstack/react-start";
import {
  getCookie,
  getWebRequest,
  setResponseStatus,
} from "@tanstack/react-start/server";

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
      console.debug("API Problem:", error);
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
