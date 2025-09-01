import type { Logger } from '@archesai/core'

import { retry } from '@archesai/core'

export const createKeyframesService = (logger: Logger) => {
  const addAbsArrayElements = (
    a: Float32Array,
    b: Float32Array
  ): Float32Array => {
    return a.map((e, i) => Math.abs(e) + Math.abs(b[i] ?? 0))
  }

  const getString = (arr: number[], isTranslation: boolean) => {
    let string = ''
    arr.forEach((value, index) => {
      let sample = value
      if (sample > 1.01 && isTranslation) {
        sample += 6
      }
      string += `${index.toString()}: (${sample.toFixed(2)})`
      if (index < arr.length - 1) {
        string += ', '
      }
    })
    return string
  }

  return {
    async getKeyframes(
      url: string,
      framerate: number,
      fn: string,
      isTranslation: boolean
    ): Promise<string> {
      logger.debug('calculating keyframes', {
        fn,
        framerate,
        isTranslation,
        url
      })

      const arrayBuffer: ArrayBuffer = await retry(
        logger,
        async () => {
          const response = await fetch(url)
          if (!response.ok) {
            throw new Error(
              `Failed to fetch file from URL: ${response.statusText}`
            )
          }
          const buffer = Buffer.from(await response.arrayBuffer())
          return buffer.buffer
        },
        5
      )

      logger.debug('buffer', { arrayBuffer })
      const audioBuffer = {
        duration: 100,
        getChannelData: (_: number) => new Float32Array(100),
        numberOfChannels: 1
      } // FIXME

      const channels: Float32Array[] = []
      for (let i = 0; i < audioBuffer.numberOfChannels; i++) {
        channels.push(audioBuffer.getChannelData(i))
      }
      const rawData = channels
        .reduce(addAbsArrayElements)
        .map((x) => x / audioBuffer.numberOfChannels)

      const samples = Math.floor(audioBuffer.duration * framerate)
      const blockSize = Math.floor(rawData.length / samples)

      let filteredData: number[] = []
      for (let i = 0; i < samples; i++) {
        const chunk = rawData.slice(i * blockSize, (i + 1) * blockSize - 1)
        const sum = chunk.reduce((a, b) => a + b, 0)
        filteredData.push(sum / chunk.length)
      }
      const max = Math.max(...filteredData)

      filteredData = filteredData
        .map((x) => x / max)
        .map(
          (x, ind) => x + ind // FIXME
        )

      return getString(filteredData, isTranslation)
    }
  }
}

export type KeyframesService = ReturnType<typeof createKeyframesService>
