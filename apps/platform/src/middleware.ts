import type { NextRequest } from 'next/server'

const auth =
  (handler: (req: NextRequest) => Promise<void>) =>
  async (req: NextRequest) => {
    // Do something here
    return handler(req)
  }

// // Or like this if you need to do something here.
export default auth(async (_req) => {
  // console.log(req) //  { session: { user: { ... } } }
  return Promise.resolve()
})

// // Read more: https://nextjs.org/docs/app/building-your-application/routing/middleware#matcher
export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)']
}
