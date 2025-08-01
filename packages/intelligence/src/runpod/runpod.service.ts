import type { ConfigService, Logger } from '@archesai/core'

import { delay, retry } from '@archesai/core'
import { IdParamsSchema, RunpodResponseSchema } from '@archesai/schemas'

export const createRunpodService = (
  configService: ConfigService,
  logger: Logger
) => {
  const monitorJob = async (
    podId: string,
    runpodJobId: string
  ): Promise<string> => {
    const token = configService.get('intelligence.runpod.token')
    if (configService.get('intelligence.runpod.mode') === 'enabled' || !token) {
      throw new Error('Runpod is not enabled due to configuration')
    }
    while (true) {
      await delay(5000)

      const response = await retry(
        logger,
        async () =>
          fetch(`https://api.runpod.ai/v2/${podId}/status/${runpodJobId}`, {
            headers: {
              Authorization: `Bearer ${token}`
            }
          }),
        5
      )

      logger.debug('runpod job status', { response })

      // fixme: add better error handling
      const responseData = await response.json()
      const parsedResponse = RunpodResponseSchema.parse(responseData)

      if (parsedResponse.status === 'COMPLETED') {
        return parsedResponse.output
      } else if (parsedResponse.status === 'FAILED') {
        throw new Error('Job failed')
      }
    }
  }

  const startJob = async (
    podId: string,
    input: Record<string, unknown>
  ): Promise<string> => {
    const token = configService.get('intelligence.runpod.token')
    if (
      configService.get('intelligence.runpod.mode') === 'disabled' ||
      !token
    ) {
      throw new Error('Runpod is not enabled due to configuration')
    }
    const response = await retry(
      logger,
      () =>
        fetch(`https://api.runpod.ai/v2/${podId}/run`, {
          body: JSON.stringify(input),
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }),
      5
    )

    logger.debug('runpod job started', { response })

    const responseData = await response.json()
    const parsedResponse = IdParamsSchema.parse(responseData)

    return parsedResponse.id
  }

  return {
    async run(
      jobId: string,
      podId: string,
      input: Record<string, unknown>
    ): Promise<string> {
      logger.debug('Starting job', { input, jobId, podId })

      const runpodJobId = await startJob(podId, input)
      return monitorJob(podId, runpodJobId)
    }
  }
}

export type RunpodService = ReturnType<typeof createRunpodService>
