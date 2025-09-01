import { z } from 'zod'

export const IntelligenceConfigSchema: z.ZodDefault<
  z.ZodOptional<
    z.ZodObject<{
      embedding: z.ZodDefault<
        z.ZodOptional<
          z.ZodObject<{
            type: z.ZodDefault<
              z.ZodEnum<{
                ollama: 'ollama'
                openai: 'openai'
              }>
            >
          }>
        >
      >
      llm: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                endpoint: z.ZodDefault<z.ZodOptional<z.ZodString>>
                token: z.ZodOptional<z.ZodString>
                type: z.ZodLiteral<'ollama'>
              }>,
              z.ZodObject<{
                endpoint: z.ZodDefault<z.ZodOptional<z.ZodString>>
                token: z.ZodOptional<z.ZodString>
                type: z.ZodLiteral<'openai'>
              }>
            ]
          >
        >
      >
      runpod: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                mode: z.ZodLiteral<'disabled'>
              }>,
              z.ZodObject<{
                mode: z.ZodLiteral<'enabled'>
                token: z.ZodOptional<z.ZodString>
              }>
            ]
          >
        >
      >
      scraper: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                mode: z.ZodLiteral<'disabled'>
              }>,
              z.ZodObject<{
                endpoint: z.ZodString
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
      speech: z.ZodDefault<
        z.ZodOptional<
          z.ZodDiscriminatedUnion<
            [
              z.ZodObject<{
                mode: z.ZodLiteral<'disabled'>
              }>,
              z.ZodObject<{
                mode: z.ZodLiteral<'enabled'>
                token: z.ZodString
              }>
            ]
          >
        >
      >
      unstructured: z.ZodDefault<
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
    }>
  >
> = z
  .object({
    embedding: z
      .object({
        type: z
          .enum(['openai', 'ollama'])
          .default('ollama')
          .describe('The embedding provider to use for vector embeddings')
      })
      .optional()
      .default({ type: 'ollama' })
      .describe('Configuration for text embedding generation'),

    llm: z
      .discriminatedUnion('type', [
        z.object({
          endpoint: z
            .string()
            .optional()
            .default('http://localhost:11434')
            .describe('Ollama server endpoint URL'),
          token: z
            .string()
            .optional()
            .describe('Optional authentication token for Ollama'),
          type: z.literal('ollama')
        }),
        z.object({
          endpoint: z
            .string()
            .optional()
            .default('https://api.openai.com/v1')
            .describe('OpenAI API endpoint (defaults to official API)'),
          token: z
            .string()
            .optional()
            .describe('OpenAI API key for authentication'),
          type: z.literal('openai')
        })
      ])
      .optional()
      .default({
        endpoint: 'http://localhost:11434',
        type: 'ollama'
      })
      .describe('Large Language Model configuration for AI processing'),

    runpod: z
      .discriminatedUnion('mode', [
        z.object({
          mode: z.literal('disabled')
        }),
        z.object({
          mode: z.literal('enabled'),
          token: z
            .string()
            .min(1)
            .optional()
            .describe('RunPod API token for serverless GPU access')
        })
      ])
      .optional()
      .default({ mode: 'disabled' })
      .describe('RunPod serverless GPU configuration for AI workloads'),

    scraper: z
      .discriminatedUnion('mode', [
        z.object({
          mode: z.literal('disabled')
        }),
        z.object({
          endpoint: z.string().describe('Web scraper service endpoint URL'),
          mode: z.literal('enabled')
        }),
        z.object({
          image: z.object({
            pullPolicy: z
              .enum(['Always', 'IfNotPresent', 'Never'])
              .optional()
              .default('IfNotPresent'),
            repository: z.string().optional().default('arches/scraper'),
            tag: z.string().optional().default('latest')
          }),
          mode: z.literal('managed'),
          resources: z.object({
            limits: z.object({
              cpu: z.string().optional().default('500m'),
              memory: z.string().optional().default('1Gi')
            }),
            requests: z.object({
              cpu: z.string().optional().default('250m'),
              memory: z.string().optional().default('512Mi')
            })
          })
        })
      ])
      .optional()
      .default({
        image: {
          pullPolicy: 'IfNotPresent',
          repository: 'arches/scraper',
          tag: 'latest'
        },
        mode: 'managed',
        resources: {
          limits: { cpu: '500m', memory: '1Gi' },
          requests: { cpu: '250m', memory: '512Mi' }
        }
      })
      .describe('Web scraping service for content extraction'),

    speech: z
      .discriminatedUnion('mode', [
        z.object({
          mode: z.literal('disabled')
        }),
        z.object({
          mode: z.literal('enabled'),
          token: z.string().describe('Speech-to-text service API token')
        })
      ])
      .optional()
      .default({ mode: 'disabled' })
      .describe('Speech recognition and text-to-speech services'),

    unstructured: z
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
            repository: z
              .string()
              .optional()
              .default(
                'downloads.unstructured.io/unstructured-io/unstructured-api'
              ),
            tag: z.string().optional().default('latest')
          }),
          mode: z.literal('managed'),
          resources: z.object({
            limits: z.object({
              cpu: z.string().optional().default('1000m'),
              memory: z.string().optional().default('2Gi')
            }),
            requests: z.object({
              cpu: z.string().optional().default('500m'),
              memory: z.string().optional().default('1Gi')
            })
          })
        })
      ])
      .optional()
      .default({
        image: {
          pullPolicy: 'IfNotPresent',
          repository:
            'downloads.unstructured.io/unstructured-io/unstructured-api',
          tag: 'latest'
        },
        mode: 'managed',
        resources: {
          limits: { cpu: '1000m', memory: '2Gi' },
          requests: { cpu: '500m', memory: '1Gi' }
        }
      })
      .describe('Unstructured.io service for document parsing and extraction')
  })
  .optional()
  .default({
    embedding: { type: 'ollama' },
    llm: {
      endpoint: 'http://localhost:11434',
      type: 'ollama'
    },
    runpod: { mode: 'disabled' },
    scraper: {
      image: {
        pullPolicy: 'IfNotPresent',
        repository: 'arches/scraper',
        tag: 'latest'
      },
      mode: 'managed',
      resources: {
        limits: { cpu: '500m', memory: '1Gi' },
        requests: { cpu: '250m', memory: '512Mi' }
      }
    },
    speech: { mode: 'disabled' },
    unstructured: {
      image: {
        pullPolicy: 'IfNotPresent',
        repository:
          'downloads.unstructured.io/unstructured-io/unstructured-api',
        tag: 'latest'
      },
      mode: 'managed',
      resources: {
        limits: { cpu: '1000m', memory: '2Gi' },
        requests: { cpu: '500m', memory: '1Gi' }
      }
    }
  })
  .describe(
    'Intelligence configuration for AI capabilities including LLMs, embeddings, web scraping, speech processing, and unstructured data handling.'
  )

export type IntelligenceConfig = z.infer<typeof IntelligenceConfigSchema>
