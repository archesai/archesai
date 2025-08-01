import { z } from 'zod'

export const createManagedFieldsSchema = ({
  image,
  pullPolicy,
  resources,
  tag
}: {
  image: string
  pullPolicy: 'Always' | 'IfNotPresent' | 'Never'
  resources: {
    limits: {
      cpu: string
      memory: string
    }
    requests: {
      cpu: string
      memory: string
    }
  }
  tag: string
}): z.ZodObject<{
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
}> => {
  return z.object({
    image: z.object({
      pullPolicy: z
        .enum(['Always', 'IfNotPresent', 'Never'])
        .optional()
        .default(pullPolicy),
      repository: z.string().optional().default(image),
      tag: z.string().optional().default(tag)
    }),
    resources: z.object({
      limits: z.object({
        cpu: z.string().optional().default(resources.limits.cpu),
        memory: z.string().optional().default(resources.limits.memory)
      }),
      requests: z.object({
        cpu: z.string().optional().default(resources.requests.cpu),
        memory: z.string().optional().default(resources.requests.memory)
      })
    })
  })
}
