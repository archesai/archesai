import type { StaticDecode } from '@sinclair/typebox'

import type {
  ArchesApiRequest,
  ArchesApiResponse,
  Controller,
  HttpInstance
} from '@archesai/core'

import {
  ArchesApiNoContentResponseSchema,
  ArchesApiUnauthorizedResponseSchema,
  IS_CONTROLLER
} from '@archesai/core'
import { AccountEntitySchema } from '@archesai/domain'

import type { AccessTokensService } from '#access-tokens/access-tokens.service'
import type { AccountsService } from '#accounts/accounts.service'
import type { AuthenticationService } from '#auth/auth.service'

import { CreateAccountRequestSchema } from '#accounts/dto/create-account.req.dto'

/**
 * Controller for managing authentication.
 */
export class AuthenticationController implements Controller {
  public readonly [IS_CONTROLLER] = true
  private readonly accessTokensService: AccessTokensService
  private readonly accountsService: AccountsService
  private readonly authenticationService: AuthenticationService

  constructor(
    accessTokensService: AccessTokensService,
    accountsService: AccountsService,
    authenticationService: AuthenticationService
  ) {
    this.accessTokensService = accessTokensService
    this.accountsService = accountsService
    this.authenticationService = authenticationService
  }

  // @UseGuards(LocalAuthGuard)
  public login(): void {
    // The LocalAuthGuard will handle the login
  }

  public async logout(
    request: ArchesApiRequest,
    reply: ArchesApiResponse
  ): Promise<void> {
    await this.authenticationService.logout(request, reply)
  }

  public async refresh(request: ArchesApiRequest): Promise<void> {
    const refreshToken: unknown = request.cookies['archesai.refreshToken']
    if (!refreshToken || typeof refreshToken !== 'string') {
      throw new Error('No refresh token provided')
    }
    const accessTokens = await this.accessTokensService.refresh(refreshToken)
    const { sub } = this.accessTokensService.verify(accessTokens.accessToken)
    await this.authenticationService.login(sub)
  }

  public async register(
    request: ArchesApiRequest & {
      body: StaticDecode<typeof CreateAccountRequestSchema>
    }
  ): Promise<void> {
    await this.accountsService.create({
      authType: 'email',
      hashed_password: request.body.password,
      provider: 'LOCAL',
      providerAccountId: request.body.email,
      userId: request.body.email
    })
  }

  public registerRoutes(app: HttpInstance) {
    app.post(
      `/auth/register`,
      {
        schema: {
          body: CreateAccountRequestSchema,
          description: `This endpoint will register you with your e-mail and password`,
          operationId: 'register',
          response: {
            204: ArchesApiNoContentResponseSchema,
            401: ArchesApiUnauthorizedResponseSchema
          },
          summary: `Register`,
          tags: ['Authentication']
        }
      },
      this.register.bind(this)
    )

    app.post(
      `/auth/login`,
      {
        schema: {
          body: CreateAccountRequestSchema,
          description: `This endpoint will log you in with your e-mail and password`,
          operationId: 'login',
          response: {
            204: ArchesApiNoContentResponseSchema,
            401: ArchesApiUnauthorizedResponseSchema
          },
          summary: `Login`,
          tags: ['Authentication']
        }
      },
      this.login.bind(this)
    )

    app.post(
      `/auth/logout`,
      {
        schema: {
          description: `This endpoint will log you out of the current session`,
          operationId: 'logout',
          response: {
            204: ArchesApiNoContentResponseSchema,
            401: ArchesApiUnauthorizedResponseSchema
          },
          summary: `Logout`,
          tags: ['Authentication']
        }
      },
      this.logout.bind(this)
    )

    app.post(
      `/auth/refresh`,
      {
        schema: {
          description: `This endpoint will refresh your access token`,
          operationId: 'refresh',
          response: {
            204: ArchesApiNoContentResponseSchema,
            401: ArchesApiUnauthorizedResponseSchema
          },
          summary: `Refresh`,
          tags: ['Authentication']
        }
      },
      this.refresh.bind(this)
    )

    app.addSchema(AccountEntitySchema)
  }
}
