import { ElevenLabs, ElevenLabsClient } from '@elevenlabs/elevenlabs-js'

import type { ConfigService } from '@archesai/core'

import { streamToBuffer } from '@archesai/core'

export const createSpeechService = (configService: ConfigService) => {
  const elevenLabsClient = new ElevenLabsClient({
    apiKey: configService.get('intelligence.speech.token')
  })

  return {
    async generate(text: string): Promise<Buffer> {
      const res = await elevenLabsClient.textToSpeech.convert(
        'pMsXgVXv3BLzUgSXRplE',
        {
          outputFormat: ElevenLabs.OutputFormat.Mp32205032,
          text: text,
          voiceSettings: {
            similarityBoost: 0.3,
            stability: 0.1,
            style: 0.2
          }
        }
      )

      return streamToBuffer(res)
    }
  }
}

export type SpeechService = ReturnType<typeof createSpeechService>
