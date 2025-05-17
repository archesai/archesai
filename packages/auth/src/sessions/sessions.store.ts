import type { SessionStore } from '@fastify/session'
import type { Session } from 'fastify'
import type { RedisClientType } from 'redis'

const noop = (_err?: unknown, _data?: Session) => {
  // no-op
}

interface NormalizedRedisClient {
  del(key: string[]): Promise<number>
  expire(key: string, ttl: number): Promise<boolean | number>
  get(key: string): Promise<null | string>
  mget(key: string[]): Promise<(null | string)[]>
  scanIterator(match: string, count: number): AsyncIterable<string>
  set(key: string, value: string, ttl?: number): Promise<null | string>
}

interface RedisStoreOptions {
  client: RedisClientType
  disableTouch?: boolean
  disableTTL?: boolean
  prefix?: string
  scanCount?: number
  serializer?: Serializer
  ttl?: ((sess: Session) => number) | number
}

interface Serializer {
  parse(s: string): Promise<Session> | Session
  stringify(s: Session): string
}

export class RedisStore implements SessionStore {
  public client: NormalizedRedisClient
  public disableTouch: boolean
  public disableTTL: boolean
  public prefix: string
  public scanCount: number
  public serializer: Serializer
  public ttl: ((sess: Session) => number) | number

  constructor(opts: RedisStoreOptions) {
    this.client = this.normalizeClient(opts.client)
    this.disableTouch = opts.disableTouch ?? false
    this.disableTTL = opts.disableTTL ?? false
    this.prefix = opts.prefix ?? 'sess:'
    this.scanCount = opts.scanCount ?? 1000
    this.serializer = opts.serializer ?? {
      parse: JSON.parse,
      stringify: JSON.stringify
    }
    this.ttl = opts.ttl ?? 86400 // default to 1 day
  }

  public async all(cb = noop) {
    const len = this.prefix.length
    try {
      const keys = await this._getAllKeys()
      if (keys.length === 0) {
        cb(null, [])
        return
      }

      const data = await this.client.mget(keys)
      const results = data.reduce<Session[]>((acc, raw, idx) => {
        if (!raw) {
          return acc
        }
        const sess = this.serializer.parse(raw) as Session
        sess.id = keys[idx]!.substring(len)
        acc.push(sess)
        return acc
      }, [])
      cb(null, results)
      return
    } catch (err) {
      cb(err)
      return
    }
  }

  public async clear(cb = noop) {
    try {
      const keys = await this._getAllKeys()
      if (!keys.length) {
        cb()
        return
      }
      await this.client.del(keys)
      cb()
      return
    } catch (err) {
      cb(err)
      return
    }
  }

  public destroy(sid: string, cb = noop) {
    const key = this.prefix + sid
    this.client
      .del([key])
      .then(() => {
        cb()
      })
      .catch((err: unknown) => {
        cb(err)
      })
  }

  public get(sid: string, cb = noop) {
    const key = this.prefix + sid
    this.client
      .get(key)
      .then((data) => {
        if (!data) {
          cb()
          return
        }
        cb(null, this.serializer.parse(data))
      })
      .catch((err: unknown) => {
        cb(err)
      })
  }

  public async ids(cb = noop) {
    const len = this.prefix.length
    try {
      const keys = await this._getAllKeys()
      cb(
        null,
        keys.map((k) => k.substring(len))
      )
      return
    } catch (err) {
      cb(err)
      return
    }
  }

  public set(sid: string, sess: Session, cb = noop) {
    const key = this.prefix + sid
    const ttl = this._getTTL(sess)
    try {
      if (ttl > 0) {
        const val = this.serializer.stringify(sess)
        if (this.disableTTL) {
          this.client
            .set(key, val)
            .then(() => {
              cb()
            })
            .catch(cb)
        } else {
          this.client
            .set(key, val, ttl)
            .then(() => {
              cb()
            })
            .catch(cb)
        }
      } else {
        this.destroy(sid, cb)
        return
      }
    } catch (err) {
      cb(err)
      return
    }
  }

  public async touch(sid: string, sess: Session, cb = noop) {
    const key = this.prefix + sid
    if (this.disableTouch || this.disableTTL) {
      cb()
      return
    }
    try {
      await this.client.expire(key, this._getTTL(sess))
      cb()
      return
    } catch (err) {
      cb(err)
      return
    }
  }

  private async _getAllKeys() {
    const pattern = this.prefix + '*'
    const keys = []
    for await (const key of this.client.scanIterator(pattern, this.scanCount)) {
      keys.push(key)
    }
    return keys
  }

  private _getTTL(sess: Session & { cookie?: { expires?: string } }) {
    if (typeof this.ttl === 'function') {
      return this.ttl(sess)
    }

    let ttl
    if (sess.cookie?.expires) {
      const ms = Number(new Date(sess.cookie.expires)) - Date.now()
      ttl = Math.ceil(ms / 1000)
    } else {
      ttl = this.ttl
    }
    return ttl
  }

  private normalizeClient(client: RedisClientType): NormalizedRedisClient {
    return {
      del: (key) => client.del(key),
      expire: (key, ttl) => client.expire(key, ttl),
      get: (key) => client.get(key),
      mget: (keys) => client.mGet(keys),
      scanIterator: (match, count) => {
        return client.scanIterator({ COUNT: count, MATCH: match })
      },
      set: (key, val, ttl) => {
        if (ttl) {
          return client.set(key, val, { EX: ttl })
        }
        return client.set(key, val)
      }
    }
  }
}
