import type { ConfigService, Logger } from '@archesai/core'

import { delay, retry } from '@archesai/core'
import { Type, Value } from '@archesai/schemas'

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

      const isValidResponse = Value.Check(
        Type.Object({
          id: Type.String(),
          output: Type.String(),
          status: Type.Union([
            Type.Literal('IN_PROGRESS'),
            Type.Literal('COMPLETED'),
            Type.Literal('FAILED')
          ])
        }),
        response.body
      )

      if (!isValidResponse) {
        throw new Error('Invalid response data')
      }

      if (response.body.status === 'COMPLETED') {
        return response.body.output
      } else if (response.body.status === 'FAILED') {
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

    const isValidResponse = Value.Check(
      Type.Object({ id: Type.String() }),
      response
    )

    if (!isValidResponse) {
      throw new Error('Invalid response data')
    }

    return response.id
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
