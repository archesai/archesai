import { z } from 'zod'

export const MonitoringConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodObject<{
      grafana: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                mode: z.ZodLiteral<'disabled'>
              }>,
              z.ZodObject<{
                mode: z.ZodLiteral<'enabled'>
              }>,
              z.ZodObject<{
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
              }>
            ]
          >
        >
      >
      loki: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                mode: z.ZodLiteral<'disabled'>
              }>,
              z.ZodObject<{
                host: z.ZodString
                mode: z.ZodLiteral<'enabled'>
              }>,
              z.ZodObject<{
                host: z.ZodDefault<z.ZodOptional<z.ZodString>>
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
              }>
            ]
          >
        >
      >
    }>
  >
> = z
  .object({
    grafana: z
      .discriminatedUnion('mode', [
        z.object({
          mode: z.literal('disabled')
        }),
        z.object({
          mode: z.literal('enabled')
        }),
        z.object({
          image: z.object({
            pullPolicy: z
              .enum(['Always', 'IfNotPresent', 'Never'])
              .optional()
              .default('IfNotPresent'),
            repository: z.string().optional().default('grafana/grafana'),
            tag: z.string().optional().default('latest')
          }),
          mode: z.literal('managed'),
          resources: z.object({
            limits: z.object({
              cpu: z.string().optional().default('200m'),
              memory: z.string().optional().default('256Mi')
            }),
            requests: z.object({
              cpu: z.string().optional().default('100m'),
              memory: z.string().optional().default('128Mi')
            })
          })
        })
      ])
      .optional()
      .default({
        image: {
          pullPolicy: 'IfNotPresent',
          repository: 'grafana/grafana',
          tag: 'latest'
        },
        mode: 'managed',
        resources: {
          limits: {
            cpu: '200m',
            memory: '256Mi'
          },
          requests: {
            cpu: '100m',
            memory: '128Mi'
          }
        }
      })
      .describe('Grafana monitoring dashboard configuration'),

    loki: z
      .discriminatedUnion('mode', [
        z.object({
          mode: z.literal('disabled')
        }),
        z.object({
          host: z.string().describe('External Loki host URL'),
          mode: z.literal('enabled')
        }),
        z.object({
          host: z.string().optional().default('http://localhost:3100'),
          image: z.object({
            pullPolicy: z
              .enum(['Always', 'IfNotPresent', 'Never'])
              .optional()
              .default('IfNotPresent'),
            repository: z.string().optional().default('grafana/loki'),
            tag: z.string().optional().default('latest')
          }),
          mode: z.literal('managed'),
          resources: z.object({
            limits: z.object({
              cpu: z.string().optional().default('200m'),
              memory: z.string().optional().default('256Mi')
            }),
            requests: z.object({
              cpu: z.string().optional().default('100m'),
              memory: z.string().optional().default('128Mi')
            })
          })
        })
      ])
      .optional()
      .default({
        host: 'http://localhost:3100',
        image: {
          pullPolicy: 'IfNotPresent',
          repository: 'grafana/loki',
          tag: 'latest'
        },
        mode: 'managed',
        resources: {
          limits: { cpu: '200m', memory: '256Mi' },
          requests: { cpu: '100m', memory: '128Mi' }
        }
      })
      .describe('Loki log aggregation service configuration')
  })
  .optional()
  .default({
    grafana: {
      image: {
        pullPolicy: 'IfNotPresent',
        repository: 'grafana/grafana',
        tag: 'latest'
      },
      mode: 'managed',
      resources: {
        limits: { cpu: '200m', memory: '256Mi' },
        requests: { cpu: '100m', memory: '128Mi' }
      }
    },
    loki: {
      host: 'http://localhost:3100',
      image: {
        pullPolicy: 'IfNotPresent',
        repository: 'grafana/loki',
        tag: 'latest'
      },
      mode: 'managed',
      resources: {
        limits: { cpu: '200m', memory: '256Mi' },
        requests: { cpu: '100m', memory: '128Mi' }
      }
    }
  })
  .describe('Monitoring configuration for Grafana and Loki services')

export type MonitoringConfig = z.infer<typeof MonitoringConfigSchema>
