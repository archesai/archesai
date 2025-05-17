import type { ConnectionOptions } from 'bullmq'

import { Queue } from 'bullmq'

export function createQueue(name: string, connection: ConnectionOptions) {
  return new Queue(name, {
    connection
  })
}
