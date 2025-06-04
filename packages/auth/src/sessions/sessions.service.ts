import { readFileSync } from 'node:fs'
import type { FastifyCookieOptions } from '@fastify/cookie'
import type { FastifySessionOptions } from '@fastify/session'
import type { Strategy } from 'passport'
import type { RedisClientType } from 'redis'

import fastifyCookie from '@fastify/cookie'
import { Authenticator } from '@fastify/passport'
import fastifySession from '@fastify/session'
import { createClient } from 'redis'

import type {
  ArchesApiRequest,
  ArchesApiResponse,
  ConfigService,
  HttpInstance
} from '@archesai/core'

import { Logger } from '@archesai/core'

import type { SessionSerializer } from '#sessions/session.serializer'

import { RedisStore } from '#sessions/sessions.store'

/**
 * Service for setting up sessions.
 */
export class SessionsService {
  private readonly configService: ConfigService
  private readonly logger = new Logger(SessionsService.name)
  private readonly sessionSerializer: SessionSerializer
  private readonly strategies: Record<string, Strategy>

  constructor(
    configService: ConfigService,
    strategies: Record<string, Strategy>,
    sessionSerializer: SessionSerializer
  ) {
    this.configService = configService
    this.sessionSerializer = sessionSerializer
    this.strategies = strategies
  }

  // public async login(userId: string, res?: ArchesApiResponse): Promise<void> {
  //   this.logger.debug('attempting to login', { userId })
  //   const accessTokens = await this.accessTokensService.create(userId)
  //   if (res) {
  //     this.logger.debug('request was passed, setting cookies')
  //     this.setCookies(res, accessTokens)
  //   } else {
  //     this.logger.debug('request was not passed, not setting cookies')
  //   }
  // }

  public async logout(
    req?: ArchesApiRequest,
    res?: ArchesApiResponse
  ): Promise<void> {
    if (res) {
      this.logger.debug('response was passed, removing cookies')
      res.clearCookie('archesai.accessToken')
      res.clearCookie('archesai.refreshToken')
      this.logger.debug('deleted cookies')
    } else {
      this.logger.debug('response was not passed, not removing cookies')
    }
    if (req) {
      this.logger.debug('request was passed, deleting cookies')
      await req.logOut()
    } else {
      this.logger.debug('request was not passed, not delsseting cookies')
    }
  }

  public setup(app: HttpInstance): void {
    this.logger.debug('setting up sessions')
    if (this.configService.get('session.enabled')) {
      // if redis is enabled, use it for the session store
      let redisStore: RedisStore | undefined
      if (this.configService.get('redis.enabled')) {
        this.logger.debug('redis enabled - using ssredis store')
        const redisClient: RedisClientType = createClient({
          password: this.configService.get('redis.auth')!,
          url: `redis://${this.configService.get('redis.host')}:${this.configService.get('redis.port').toString()}`,
          ...(this.configService.get('redis.ca') ?
            {
              socket: {
                ca: readFileSync(this.configService.get('redis.ca')!),
                host: this.configService.get('redis.host'),
                rejectUnauthorized: false,
                tls: true
              }
            }
          : {})
        })

        redisClient.on('error', (error: unknown) => {
          this.logger.error('redis client error', { error })
        })
        redisClient.on('connect', () => {
          this.logger.debug('redis client connected')
        })
        redisClient.connect().catch(console.error)
        redisStore = new RedisStore({
          client: redisClient
        })
      }

      // Register the cookie plugin first, so that sessions can use it for signing
      app.register(fastifyCookie, {
        secret: this.configService.get('session.secret')
      } satisfies FastifyCookieOptions)

      // Register the session plugin
      app.register(fastifySession, {
        cookie: {
          httpOnly: true,
          maxAge: 24 * 60 * 60 * 1000, // e.g., 1 hour in milliseconds
          sameSite: 'lax',
          secure: this.configService.get('tls.enabled') // Set to true if using HTTPS in production
        },
        secret: this.configService.get('session.secret'),
        ...(redisStore ?
          {
            store: redisStore
          }
        : {}),
        saveUninitialized: false // Do not save session if unmodified
      } satisfies FastifySessionOptions)

      const fastifyPassport = new Authenticator()

      fastifyPassport.registerUserSerializer(
        this.sessionSerializer.serializeUser.bind(this.sessionSerializer)
      )
      fastifyPassport.registerUserDeserializer(
        this.sessionSerializer.deserializeUser.bind(this.sessionSerializer)
      )

      app.register(fastifyPassport.initialize())
      app.register(fastifyPassport.secureSession())

      // Register the strategies
      for (const [name, strategy] of Object.entries(this.strategies)) {
        fastifyPassport.use(name, strategy)
      }

      this.logger.debug('session setup completed')
    } else {
      this.logger.debug('session disabled - skipping')
    }
  }

  // private setCookies(
  //   res: ArchesApiResponse,
  //   accessTokens: AccessTokenEntity
  // ): void {
  //   res.cookie('archesai.accessToken', accessTokens.accessToken, {
  //     httpOnly: true,
  //     maxAge: 15 * 60 * 1000, // 15 minutes for access token
  //     sameSite: 'none',
  //     secure: true
  //   })
  //   res.cookie('archesai.refreshToken', accessTokens.refreshToken, {
  //     httpOnly: true,
  //     maxAge: 7 * 24 * 60 * 60 * 1000, // 7 days for refresh token
  //     sameSite: 'none',
  //     secure: true,
  //     signed: true
  //   })
  //   this.logger.debug('set cookies in response')
  // }
}
