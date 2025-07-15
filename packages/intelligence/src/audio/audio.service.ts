import { randomUUID } from 'node:crypto'
import { unlinkSync, writeFileSync } from 'node:fs'
import os from 'node:os'
import path from 'node:path'

import ffmpeg from 'fluent-ffmpeg'

import type { Logger } from '@archesai/core'

import { retry } from '@archesai/core'
import { Type, Value } from '@archesai/schemas'

export const createAudioService = (logger: Logger) => {
  return {
    async splitAudio(url: string): Promise<{
      bassSrc: string
      drumsSrc: string
    }> {
      logger.debug("Hitting moises' API...")
      const response = await retry(
        logger,
        () =>
          fetch('https://developer-api.moises.ai/api/job', {
            body: JSON.stringify({
              name: randomUUID(),
              params: {
                inputUrl: url
              },
              workflow: 'archesai-workflow'
            }),
            headers: {
              Authorization: '5fa360fa-9974-47fc-bcb8-39142bf4dcea',
              'Content-Type': 'application/json'
            }
          }),
        5
      )
      logger.debug('got response from moises', {
        response
      })
      const validatedResponse = Value.Parse(
        Type.Object({
          id: Type.String()
        }),
        response
      )

      while (true) {
        logger.debug('checking moises job', { validatedResponse })
        await new Promise((resolve) => setTimeout(resolve, 5000))
        const response = await retry(
          logger,
          () =>
            fetch(
              'https://developer-api.moises.ai/api/job/' + validatedResponse.id,
              {
                headers: {
                  Authorization: '5fa360fa-9974-47fc-bcb8-39142bf4dcea'
                }
              }
            ),
          5
        )
        logger.debug('got response from moises', {
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
          response.json()
        )

        if (moisesCheckJobResponse.status == 'SUCCEEDED') {
          const bassSrc = moisesCheckJobResponse.result.Bass
          const drumsSrc = moisesCheckJobResponse.result.Drums
          logger.debug('got bass and drums src', { bassSrc, drumsSrc })

          return { bassSrc, drumsSrc }
        } else if (moisesCheckJobResponse.status === 'FAILED') {
          throw new Error("Moises' job failed")
        }
      }
    },

    async trimAudio(
      url: string,
      startTime: number,
      duration: number
    ): Promise<string> {
      const inputTmpPath = path.join(os.tmpdir(), 'original.mp3')
      const outputTmpPath = path.join(os.tmpdir(), 'trimmed.mp3')
      const response = await (await fetch(url)).arrayBuffer()
      writeFileSync(inputTmpPath, Buffer.from(response))

      return new Promise<string>((resolve, reject) => {
        ffmpeg.ffprobe(inputTmpPath, (err: unknown, data) => {
          if (err) {
            logger.error('error getting audio duration', { err })
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
}

export type AudioService = ReturnType<typeof createAudioService>
