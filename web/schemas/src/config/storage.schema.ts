import { z } from 'zod'

const BaseStorageConfigSchema = z.object({
  accesskey: z
    .string()
    .optional()
    .default('minioadmin')
    .describe('MinIO/S3 access key ID for authentication'),
  bucket: z
    .string()
    .optional()
    .default('archesai')
    .describe('S3 bucket name for file storage'),
  endpoint: z
    .string()
    .optional()
    .default('http://localhost:9000')
    .describe('MinIO server endpoint URL'),
  secretkey: z
    .string()
    .optional()
    .default('minioadmin')
    .describe('MinIO/S3 secret access key for authentication')
})

export const StorageConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodDiscriminatedUnion<
      [
        z.ZodObject<{
          mode: z.ZodLiteral<'disabled'>
        }>,
        z.ZodObject<{
          accesskey: z.ZodDefault<z.ZodOptional<z.ZodString>>
          bucket: z.ZodDefault<z.ZodOptional<z.ZodString>>
          endpoint: z.ZodDefault<z.ZodOptional<z.ZodString>>
          mode: z.ZodLiteral<'enabled'>
          secretkey: z.ZodDefault<z.ZodOptional<z.ZodString>>
        }>,
        z.ZodObject<{
          accesskey: z.ZodDefault<z.ZodOptional<z.ZodString>>
          bucket: z.ZodDefault<z.ZodOptional<z.ZodString>>
          endpoint: z.ZodDefault<z.ZodOptional<z.ZodString>>
          image: z.ZodObject<{
            pullPolicy: z.ZodDefault<
              z.ZodOptional<
                z.ZodEnum<{
                  Always: 'Always'
                  IfNotPresent: 'IfNotPresent'
                  Never: 'Never'
                }>
              >
            >
            repository: z.ZodDefault<z.ZodOptional<z.ZodString>>
            tag: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
          mode: z.ZodLiteral<'managed'>
          persistence: z.ZodObject<{
            enabled: z.ZodDefault<z.ZodOptional<z.ZodBoolean>>
            size: z.ZodDefault<z.ZodOptional<z.ZodString>>
          }>
          resources: z.ZodObject<{
            limits: z.ZodObject<{
              cpu: z.ZodDefault<z.ZodOptional<z.ZodString>>
              memory: z.ZodDefault<z.ZodOptional<z.ZodString>>
            }>
            requests: z.ZodObject<{
              cpu: z.ZodDefault<z.ZodOptional<z.ZodString>>
              memory: z.ZodDefault<z.ZodOptional<z.ZodString>>
            }>
          }>
          secretkey: z.ZodDefault<z.ZodOptional<z.ZodString>>
        }>
      ]
    >
  >
> = z
  .discriminatedUnion('mode', [
    z.object({
      mode: z.literal('disabled')
    }),
    BaseStorageConfigSchema.extend({
      mode: z.literal('enabled')
    }),
    BaseStorageConfigSchema.extend({
      image: z.object({
        pullPolicy: z
          .enum(['Always', 'IfNotPresent', 'Never'])
          .optional()
          .default('IfNotPresent'),
        repository: z.string().optional().default('minio/minio'),
        tag: z.string().optional().default('latest')
      }),
      mode: z.literal('managed'),
      persistence: z.object({
        enabled: z
          .boolean()
          .optional()
          .default(true)
          .describe('Enable persistent storage for MinIO data'),
        size: z
          .string()
          .optional()
          .default('20Gi')
          .describe('Size of persistent volume for object storage')
      }),
      resources: z.object({
        limits: z.object({
          cpu: z.string().optional().default('500m'),
          memory: z.string().optional().default('512Mi')
        }),
        requests: z.object({
          cpu: z.string().optional().default('250m'),
          memory: z.string().optional().default('256Mi')
        })
      })
    })
  ])
  .optional()
  .default({
    accesskey: 'minioadmin',
    bucket: 'arches',
    endpoint: 'http://localhost:9000',
    image: {
      pullPolicy: 'IfNotPresent',
      repository: 'minio/minio',
      tag: 'latest'
    },
    mode: 'managed',
    persistence: {
      enabled: true,
      size: '20Gi'
    },
    resources: {
      limits: { cpu: '500m', memory: '512Mi' },
      requests: { cpu: '250m', memory: '256Mi' }
    },
    secretkey: 'minioadmin'
  })
  .describe('Object storage configuration for MinIO or S3-compatible services')

export type StorageConfig = z.infer<typeof StorageConfigSchema>
