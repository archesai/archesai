import type { GetSession200 } from "@archesai/client";
import { getSession } from "@archesai/client";
import { createServerFn } from "@tanstack/react-start";
import { getWebRequest } from "@tanstack/react-start/server";

const getSessionServer = createServerFn({
  method: "GET",
}).handler(async (): Promise<GetSession200 | null> => {
  const { headers } = getWebRequest();
  try {
    const result = await getSession("current", {
      credentials: "include",
      headers,
    });
    return result;
  } catch {
    /* empty */
  }
  return null;
});

export default getSessionServer;
