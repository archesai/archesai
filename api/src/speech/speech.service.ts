// elevenlabs/elevenlabs.service.ts
import { Injectable } from "@nestjs/common";
import { ElevenLabs, ElevenLabsClient } from "elevenlabs";
import internal from "stream";

@Injectable()
export class SpeechService {
  private readonly client: ElevenLabsClient;

  constructor() {
    this.client = new ElevenLabsClient({
      apiKey: "sk_d6c86e029888836fa68bda6cc3cb40de0aefb9fbb2ca76a9",
    });
  }

  async generateSpeech(text: string): Promise<Buffer> {
    const res = await this.client.textToSpeech.convert("pMsXgVXv3BLzUgSXRplE", {
      optimize_streaming_latency: ElevenLabs.OptimizeStreamingLatency.Zero,
      output_format: ElevenLabs.OutputFormat.Mp32205032,
      text: text,
      voice_settings: {
        similarity_boost: 0.3,
        stability: 0.1,
        style: 0.2,
      },
    });

    // Convert the stream to a buffer
    const audioBuffer = await this.streamToBuffer(res);

    return audioBuffer;
  }

  private async streamToBuffer(stream: internal.Readable): Promise<Buffer> {
    const chunks: Buffer[] = [];
    for await (const chunk of stream) {
      chunks.push(Buffer.isBuffer(chunk) ? chunk : Buffer.from(chunk));
    }
    return Buffer.concat(chunks);
  }
}
