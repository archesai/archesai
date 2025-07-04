import type { Controller, HttpInstance } from '@archesai/core'

import { ArchesApiNoContentResponseSchema, IS_CONTROLLER } from '@archesai/core'
import { LegacyRef } from '@archesai/domain'

/**
 * Controller for managing OAuth.
 */
export class OAuthController implements Controller {
  public readonly [IS_CONTROLLER] = true

  public firebase(): Promise<void> {
    throw new Error('Handler for POST /oauth/firebase not implemented')
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      // @UseGuards(AuthGuard('firebase-auth'))
      `/oauth/firebase`,
      {
        schema: {
          description: 'Authenticate with Firebase using OAuth',
          operationId: 'firebase',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema)
          },
          summary: 'Authenticate with Firebase',
          tags: ['OAuth']
        }
      },
      this.firebase.bind(this)
    )

    app.get(
      // @UseGuards(AuthGuard('twitter'))
      `/oauth/twitter`,
      {
        schema: {
          description: 'Redirects user to Twitter for authentication.',
          operationId: 'twitter',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema)
          },
          summary: 'Redirect to Twitter OAuth',
          tags: ['OAuth']
        }
      },
      this.twitter.bind(this)
    )

    app.get(
      // @UseGuards(AuthGuard('twitter'))
      `/oauth/twitter/callback`,
      {
        schema: {
          description:
            'Receives the Twitter OAuth callback after successful authentication.',
          operationId: 'twitterCallback',
          response: {
            204: LegacyRef(ArchesApiNoContentResponseSchema)
          },
          summary: 'Handle Twitter OAuth callback',
          tags: ['OAuth']
        }
      },
      this.twitterCallback.bind(this)
    )
  }

  public twitter(): Promise<void> {
    throw new Error('Handler for GET /oauth/twitter not implemented')
  }

  public twitterCallback(): Promise<void> {
    throw new Error('Handler for GET /oauth/twitter/callback not implemented')
  }
}
