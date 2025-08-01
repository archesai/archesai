import { z } from 'zod'

const BaseIngressConfigSchema = z.object({
  domain: z
    .string()
    .optional()
    .default('archesai.dev')
    .describe('Primary domain name for ingress routing'),
  tls: z
    .object({
      enabled: z
        .boolean()
        .optional()
        .default(true)
        .describe('Enable TLS/SSL encryption for HTTPS traffic'),
      secretName: z
        .string()
        .optional()
        .default('archesai-tls')
        .describe('Kubernetes secret name containing TLS certificate and key')
    })
    .optional()
    .default({
      enabled: true,
      secretName: 'archesai-tls'
    })
    .describe(
      'TLS configuration for ingress, including certificate management and encryption settings'
    )
})

export const IngressConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          domain: z.ZodDefault<z.ZodOptional<z.ZodString>>
          mode: z.ZodLiteral<'enabled'>
          tls: z.ZodDefault<
            z.ZodOptional<
              z.ZodObject<{
                enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
                secretName: z.ZodDefault<z.ZodOptional<z.ZodString>>
              }>
            >
          >
        }>,
        z.ZodObject<{
          domain: z.ZodDefault<z.ZodOptional<z.ZodString>>
          mode: z.ZodLiteral<'managed'>
          tls: z.ZodObject<{
            enabled: z.ZodDefault<z.ZodBoolean>
            issuer: z.ZodDefault<z.ZodString>
            secretName: z.ZodDefault<z.ZodString>
          }>
        }>
      ]
    >
  >
> = z
  .discriminatedUnion('mode', [
    z.object({
      mode: z.literal('disabled')
    }),
    BaseIngressConfigSchema.extend({
      mode: z.literal('enabled')
    }),
    BaseIngressConfigSchema.extend({
      mode: z.literal('managed'),
      tls: z.object({
        enabled: z
          .boolean()
          .default(true)
          .describe('Enable TLS/SSL encryption for HTTPS traffic'),
        issuer: z
          .string()
          .default('letsencrypt-staging')
          .describe(
            'Cert-manager ClusterIssuer name for automatic certificate generation'
          ),
        secretName: z
          .string()
          .default('archesai-tls')
          .describe(
            'Kubernetes secret name for storing generated TLS certificate'
          )
      })
    })
  ])
  .optional()
  .default({
    domain: 'archesai.dev',
    mode: 'managed',
    tls: {
      enabled: true,
      issuer: 'letsencrypt-staging',
      secretName: 'archesai-tls'
    }
  })
  .describe(
    'Ingress configuration for routing external traffic to internal services. Supports TLS/SSL encryption and automatic certificate management.'
  )

export type IngressConfig = z.infer<typeof IngressConfigSchema>
