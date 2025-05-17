import { Type } from '@sinclair/typebox'
import { Value } from '@sinclair/typebox/value'

import type { ConfigService, FetcherService } from '@archesai/core'

import { delay, Logger, retry } from '@archesai/core'

/**
 * Service for interacting with the RunPod API.
 */
export class RunpodService {
  private readonly configService: ConfigService
  private readonly fetcherService: FetcherService
  private readonly logger = new Logger(RunpodService.name)

  constructor(configService: ConfigService, fetcherService: FetcherService) {
    this.configService = configService
    this.fetcherService = fetcherService
  }

  /**
   * Executes a job on RunPod and monitors its status until completion.
   * @param jobId - The unique identifier for the job to be executed.
   * @param podId - The unique identifier for the pod where the job will run.
   * @param input - A record containing the input parameters required for the job.
   * @returns A promise that resolves to the output of the completed job as a string.
   * @throws Will throw an error if the response data from RunPod is invalid or if the job fails.
   */
  public async run(
    jobId: string,
    podId: string,
    input: Record<string, unknown>
  ): Promise<string> {
    this.logger.debug('Starting job', { input, jobId, podId })

    const runpodJobId = await this.startJob(podId, input)
    return this.monitorJob(podId, runpodJobId)
  }

  private async monitorJob(
    podId: string,
    runpodJobId: string
  ): Promise<string> {
    while (true) {
      await delay(5000)

      const response = await retry(
        this.logger,
        () =>
          this.fetcherService.get(
            `https://api.runpod.ai/v2/${podId}/status/${runpodJobId}`,
            {
              Authorization: `Bearer ${this.configService.get('runpod.token')}`
            }
          ),
        5
      )

      this.logger.debug('runpod job status', { response })

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
        response
      )

      if (!isValidResponse) {
        throw new Error('Invalid response data')
      }

      if (response.status === 'COMPLETED') {
        return response.output
      } else if (response.status === 'FAILED') {
        throw new Error('Job failed')
      }
    }
  }

  private async startJob(
    podId: string,
    input: Record<string, unknown>
  ): Promise<string> {
    const response = await retry(
      this.logger,
      () =>
        this.fetcherService.post(
          `https://api.runpod.ai/v2/${podId}/run`,
          input,
          {
            Authorization: `Bearer ${this.configService.get('runpod.token')}`
          }
        ),
      5
    )

    this.logger.debug('runpod job started', { response })

    const isValidResponse = Value.Check(
      Type.Object({ id: Type.String() }),
      response
    )

    if (!isValidResponse) {
      throw new Error('Invalid response data')
    }

    return response.id
  }
}
