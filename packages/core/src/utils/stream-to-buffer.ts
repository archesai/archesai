import { buffer } from 'node:stream/consumers'

export const streamToBuffer = async (
  stream: ReadableStream<ArrayBufferLike | Uint8Array>
): Promise<Buffer> => {
  return await buffer(stream)
}
