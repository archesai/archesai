import type { ConfigService, Logger } from '@archesai/core'

import { delay, retry } from '@archesai/core'

export const createRunpodService = (
  configService: ConfigService,
  logger: Logger
) => {
  const monitorJob = async (
    podId: string,
    runpodJobId: string
  ): Promise<string> => {
    while (true) {
      await delay(5000)

      const response = await retry(
        logger,
        async () =>
          fetch(`https://api.runpod.ai/v2/${podId}/status/${runpodJobId}`, {
            headers: {
              Authorization: `Bearer ${configService.get('runpod.token')}`
            }
          }),
        5
      )

      logger.debug('runpod job status', { response })

      // fixme: add better error handling
      const responseData = (await response.json()) as {
        id: string
        output: string
        status: 'COMPLETED' | 'FAILED' | 'IN_PROGRESS'
      }

      if (responseData.status === 'COMPLETED') {
        return responseData.output
      } else if (responseData.status === 'FAILED') {
        throw new Error('Job failed')
      }
    }
  }

  const startJob = async (
    podId: string,
    input: Record<string, unknown>
  ): Promise<string> => {
    const response = await retry(
      logger,
      () =>
        fetch(`https://api.runpod.ai/v2/${podId}/run`, {
          body: JSON.stringify(input),
          headers: {
            Authorization: `Bearer ${configService.get('runpod.token')}`,
            'Content-Type': 'application/json'
          },
          method: 'POST'
        }),
      5
    )

    logger.debug('runpod job started', { response })

    const responseData = (await response.json()) as { id: string }

    return responseData.id
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
