import { randomUUID } from 'node:crypto'
import { unlinkSync, writeFileSync } from 'node:fs'
import os from 'node:os'
import path from 'node:path'

import ffmpeg from 'fluent-ffmpeg'

import type { FetcherService } from '@archesai/core'

import { Logger, retry } from '@archesai/core'
import { Type, Value } from '@archesai/schemas'

/**
 * Service for processing audio files.
 */
export class AudioService {
  private readonly fetcherService: FetcherService
  private readonly logger = new Logger(AudioService.name)

  constructor(fetcherService: FetcherService) {
    this.fetcherService = fetcherService
  }

  /**
   * Splits an audio file into its components (e.g., bass and drums) using Moises' API.
   * @param url - The URL of the audio file to be processed.
   * @returns A promise that resolves to an object containing the URLs of the bass and drums audio components.
   * @throws Will throw an error if the Moises job fails or if the API responses do not match the expected structure.
   * This method interacts with Moises' API to process the audio file. It first submits a job to the API
   * and then continuously polls the job status until it succeeds or fails. The method uses a retry mechanism
   * for API calls to handle transient errors.
   * @example
   * ```typescript
   * const { bassSrc, drumsSrc } = await audioService.splitAudio('https://example.com/audio.mp3');
   * console.log('Bass URL:', bassSrc);
   * console.log('Drums URL:', drumsSrc);
   * ```
   */
  public async splitAudio(url: string) {
    this.logger.debug("Hitting moises' API...")
    const response = await retry(
      this.logger,
      () =>
        this.fetcherService.post(
          'https://developer-api.moises.ai/api/job',
          {
            name: randomUUID(),
            params: {
              inputUrl: url
            },
            workflow: 'archesai-workflow'
          },
          {
            Authorization: '5fa360fa-9974-47fc-bcb8-39142bf4dcea',
            'Content-Type': 'application/json'
          }
        ),
      5
    )
    this.logger.debug('got response from moises', {
      response
    })
    const validatedResponse = Value.Parse(
      Type.Object({
        id: Type.String()
      }),
      response
    )

    while (true) {
      this.logger.debug('checking moises job', { validatedResponse })
      await new Promise((resolve) => setTimeout(resolve, 5000))
      const response = await retry<{
        data: Record<string, unknown>
      }>(
        this.logger,
        () =>
          this.fetcherService.get(
            'https://developer-api.moises.ai/api/job/' + validatedResponse.id,
            {
              Authorization: '5fa360fa-9974-47fc-bcb8-39142bf4dcea'
            }
          ),
        5
      )
      this.logger.debug('got response from moises', {
        response
      })
      const moisesCheckJobResponse = Value.Parse(
        Type.Object({
          result: Type.Object({
            Bass: Type.String(),
            Drums: Type.String()
          }),
          status: Type.String()
        }),
        response.data
      )

      if (moisesCheckJobResponse.status == 'SUCCEEDED') {
        const bassSrc = moisesCheckJobResponse.result.Bass
        const drumsSrc = moisesCheckJobResponse.result.Drums
        this.logger.debug('got bass and drums src', { bassSrc, drumsSrc })

        return { bassSrc, drumsSrc }
      } else if (moisesCheckJobResponse.status === 'FAILED') {
        throw new Error("Moises' job failed")
      }
    }
  }

  /**
   * Trims an audio file from a given URL to a specified start time and duration.
   * @param url - The URL of the audio file to be trimmed.
   * @param startTime - The start time (in seconds) from which the audio should be trimmed.
   * @param duration - The duration (in seconds) of the trimmed audio segment.
   * @returns A promise that resolves to the URL or path of the trimmed audio file.
   * @throws Will throw an error if there is an issue fetching the audio file,
   *         determining its duration, or processing the trimming operation.
   */
  public async trimAudio(
    url: string,
    startTime: number,
    duration: number
  ): Promise<string> {
    const inputTmpPath = path.join(os.tmpdir(), 'original.mp3')
    const outputTmpPath = path.join(os.tmpdir(), 'trimmed.mp3')
    const response = await this.fetcherService.get<Buffer>(url)
    writeFileSync(inputTmpPath, response)

    return new Promise<string>((resolve, reject) => {
      ffmpeg.ffprobe(inputTmpPath, (err: unknown, data) => {
        if (err) {
          this.logger.error('error getting audio duration', { err })
          reject(new Error('error getting audio duration'))
        }
        if ((data.format.duration ?? 0) < startTime + duration) {
          unlinkSync(inputTmpPath)
          resolve(url)
        }
      })

      ffmpeg(inputTmpPath)
        .setStartTime(startTime)
        .setDuration(duration)
        .output(outputTmpPath)
        .on('end', () => {
          // this.storageService
          //   .uploadFromFile(
          //     'audio/' + new Date().valueOf().toString() + '.mp3',
          //     {
          //       buffer: readFileSync(outputTmpPath),
          //       mimetype: 'audio/mp3',
          //       originalname: 'audio.mp3'
          //     }
          //   )
          //   .then((file) => {
          //     unlinkSync(outputTmpPath)
          //     unlinkSync(inputTmpPath)
          //     resolve(file.read || '')
          //   })
          //   .catch((err: unknown) => {
          //     this.logger.error({ err }, 'error uploading file')
          //     reject(new Error('error uploading file'))
          //   })
        })
        .on('error', reject)
        .run()
    })
  }
}
