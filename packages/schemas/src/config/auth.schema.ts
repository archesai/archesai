import { z } from 'zod'

const BaseAuthConfigSchema = z.object({
  firebase: z
    .discriminatedUnion('mode', [
      z.object({
        mode: z.literal('disabled')
      }),
      z.object({
        clientEmail: z
          .string()
          .describe('Firebase service account client email address'),
        mode: z.literal('enabled'),
        privateKey: z
          .string()
          .describe('Firebase service account private key (PEM format)'),
        projectId: z.string().describe('Firebase project ID for authentication')
      })
    ])
    .optional()
    .default({ mode: 'disabled' })
    .describe(
      'Firebase authentication configuration. Enables Google Firebase Auth integration for user authentication and authorization.'
    ),
  local: z
    .discriminatedUnion('mode', [
      z.object({
        mode: z.literal('disabled')
      }),
      z.object({
        mode: z.literal('enabled')
      })
    ])
    .optional()
    .default({ mode: 'enabled' })
    .describe(
      'Local authentication configuration. Provides username/password authentication stored in your database.'
    ),
  twitter: z
    .discriminatedUnion('mode', [
      z.object({
        mode: z.literal('disabled')
      }),
      z.object({
        callbackURL: z
          .string()
          .describe(
            'OAuth callback URL that Twitter will redirect to after authentication'
          ),
        consumerKey: z
          .string()
          .describe('Twitter API consumer key (API key) from your Twitter app'),
        consumerSecret: z
          .string()
          .describe(
            'Twitter API consumer secret (API secret key) from your Twitter app'
          ),
        mode: z.literal('enabled')
      })
    ])
    .optional()
    .default({ mode: 'disabled' })
    .describe(
      'Twitter OAuth authentication configuration. Enables "Sign in with Twitter" functionality for users.'
    )
})

export const AuthConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          firebase: z.ZodDefault<
            z.ZodOptional<
              z.ZodDiscriminatedUnion<
                [
                  z.ZodObject<{
                    mode: z.ZodLiteral<'disabled'>
                  }>,
                  z.ZodObject<{
                    clientEmail: z.ZodString
                    mode: z.ZodLiteral<'enabled'>
                    privateKey: z.ZodString
                    projectId: z.ZodString
                  }>
                ]
              >
            >
          >
          local: z.ZodDefault<
            z.ZodOptional<
              z.ZodDiscriminatedUnion<
                [
                  z.ZodObject<{
                    mode: z.ZodLiteral<'disabled'>
                  }>,
                  z.ZodObject<{
                    mode: z.ZodLiteral<'enabled'>
                  }>
                ]
              >
            >
          >
          mode: z.ZodLiteral<'enabled'>
          twitter: z.ZodDefault<
            z.ZodOptional<
              z.ZodDiscriminatedUnion<
                [
                  z.ZodObject<{
                    mode: z.ZodLiteral<'disabled'>
                  }>,
                  z.ZodObject<{
                    callbackURL: z.ZodString
                    consumerKey: z.ZodString
                    consumerSecret: z.ZodString
                    mode: z.ZodLiteral<'enabled'>
                  }>
                ]
              >
            >
          >
        }>
      ]
    >
  >
> = z
  .discriminatedUnion('mode', [
    z.object({
      mode: z.literal('disabled')
    }),
    BaseAuthConfigSchema.extend({
      mode: z.literal('enabled')
    })
  ])
  .optional()
  .default({
    firebase: {
      mode: 'disabled'
    },
    local: { mode: 'enabled' },
    mode: 'enabled',
    twitter: { mode: 'disabled' }
  })
  .describe(
    'Authentication configuration for the API server. This includes Firebase, local, and Twitter authentication options. Each option can be enabled or disabled independently. The default mode is "enabled" with local authentication enabled.'
  )

export type AuthConfig = z.infer<typeof AuthConfigSchema>
