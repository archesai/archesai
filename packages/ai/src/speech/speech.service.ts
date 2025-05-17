import { ElevenLabs, ElevenLabsClient } from 'elevenlabs'

import type { ConfigService } from '@archesai/core'

import { streamToBuffer } from '@archesai/core'

/**
 * Service for generating speech from text.
 */
export class SpeechService {
  private readonly configService: ConfigService
  private readonly elevenLabs: ElevenLabsClient

  constructor(configService: ConfigService) {
    this.configService = configService
    this.elevenLabs = new ElevenLabsClient({
      apiKey: this.configService.get('speech.token')
    })
  }

  /**
   * Generates speech from the provided text.
   * @param text - The text to convert to speech.
   * @returns A promise that resolves with the generated speech as a buffer.
   */
  public async generate(text: string): Promise<Buffer> {
    const res = await this.elevenLabs.textToSpeech.convert(
      'pMsXgVXv3BLzUgSXRplE',
      {
        output_format: ElevenLabs.OutputFormat.Mp32205032,
        text: text,
        voice_settings: {
          similarity_boost: 0.3,
          stability: 0.1,
          style: 0.2
        }
      }
    )

    return streamToBuffer(res)
  }
}
