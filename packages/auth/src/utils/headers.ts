import type { FastifyReply } from 'fastify'

export const getHeaders = (
  headers: Record<string, string | string[] | undefined>
): Headers => {
  const headersObj = new Headers()
  Object.entries(headers).forEach(([key, value]) => {
    if (value) {
      if (Array.isArray(value)) {
        value.forEach((v) => {
          headersObj.append(key, v)
        })
      } else {
        headersObj.append(key, value)
      }
    }
  })

  return headersObj
}

export const setHeaders = (headers: Headers, response: FastifyReply): void => {
  headers.forEach((value, key) => {
    response.header(key, value)
  })
}

// Optional: Add helper methods to fastify instance
// app.decorate(
//   'authHandler',
//   async (
//     req: FastifyRequest,
//     reply: FastifyReply,
//     beforeSend?: (
//       response: Response,
//       responseText: null | Record<string, unknown>
//     ) => Promise<void>
//   ) => {
//     // Reusable auth handler logic that can be called from other routes
//     const url = new URL(
//       req.url,
//       `http://${req.headers.host?.toString() ?? ''}`
//     )

//     const headers = new Headers()
//     Object.entries(req.headers).forEach(([key, value]) => {
//       if (value) headers.append(key, value.toString())
//     })

//     const formattedRequest = new Request(url.toString(), {
//       body: req.body ? JSON.stringify(req.body) : undefined,
//       headers,
//       method: req.method
//     })

//     const response = await authService.handler(formattedRequest)

//     // Get response text once
//     const responseText = response.body ? await response.text() : null

//     // Forward response to client
//     reply.status(response.status)
//     response.headers.forEach((value, key) => {
//       reply.header(key, value)
//     })

//     // Run callback if provided
//     if (beforeSend) {
//       const responseJson = (
//         responseText ?
//           JSON.parse(responseText)
//         : null) as null | Record<string, unknown>
//       await beforeSend(response, responseJson)
//     }

//     reply.send(responseText)
//   }
// )
