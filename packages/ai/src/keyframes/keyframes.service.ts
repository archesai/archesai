import { Logger, retry } from '@archesai/core'

/**
 * Service for calculating keyframes from audio data.
 */
export class KeyframesService {
  private readonly logger = new Logger(KeyframesService.name)

  /**
   * Adds the absolute values of the elements of two Float32Arrays.
   * @param a - Float32Array
   * @param b - Float32Array
   * @returns Float32Array
   */
  public addAbsArrayElements = (
    a: Float32Array,
    b: Float32Array
  ): Float32Array => {
    return a.map((e, i) => Math.abs(e) + Math.abs(b[i] ?? 0))
  }

  /**
   * Processes audio data to calculate keyframes based on the provided parameters.
   * @param url - The URL of the audio file to process.
   * @param framerate - The desired framerate for the keyframes.
   * @param fn - A mathematical function represented as a string to evaluate on the normalized data.
   * @param isTranslation - A flag indicating whether the output should be formatted for translation.
   * @returns A formatted string representing the calculated keyframes.
   */
  public async getKeyframes(
    url: string,
    framerate: number,
    fn: string,
    isTranslation: boolean
  ): Promise<string> {
    this.logger.debug('calculating keyframes', {
      fn,
      framerate,
      isTranslation,
      url
    })

    const arrayBuffer: ArrayBuffer = await retry(
      this.logger,
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

    // Decode the audio data into an AudioBuffer
    // const audioBuffer: typeof AudioBuffer =
    //   await context.decodeAudioData(arrayBuffer)
    this.logger.debug('buffer', { arrayBuffer })
    const audioBuffer = {
      duration: 100,
      getChannelData: (_: number) => new Float32Array(100),
      numberOfChannels: 1
    } // FIXME

    // Combine all audio channels into one average channel
    const channels: Float32Array[] = []
    for (let i = 0; i < audioBuffer.numberOfChannels; i++) {
      channels.push(audioBuffer.getChannelData(i))
    }
    const rawData = channels
      .reduce(this.addAbsArrayElements)
      .map((x) => x / audioBuffer.numberOfChannels)

    // Determine the number of samples and block size
    const samples = Math.floor(audioBuffer.duration * framerate)
    const blockSize = Math.floor(rawData.length / samples)

    // Filter and normalize the data
    let filteredData: number[] = []
    for (let i = 0; i < samples; i++) {
      const chunk = rawData.slice(i * blockSize, (i + 1) * blockSize - 1)
      const sum = chunk.reduce((a, b) => a + b, 0)
      filteredData.push(sum / chunk.length)
    }
    const max = Math.max(...filteredData)

    // Evaluate the mathematical function on normalized data
    filteredData = filteredData
      .map((x) => x / max)
      .map(
        (x, ind) =>
          // evaluate(
          //   fn.replace('x', x.toString()).replace('y', ind.toString())
          // ) satisfies number
          x + ind // FIXME
      )

    // Convert the data into a formatted string
    return this.getString(filteredData, isTranslation)
  }

  /**
   * Formats the keyframe data to a string.
   * @param arr - The array of keyframe data.
   * @param isTranslation - Determines if translation adjustment is needed.
   * @returns A formatted string representation of the keyframes.
   */
  public getString = (arr: number[], isTranslation: boolean) => {
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
}
