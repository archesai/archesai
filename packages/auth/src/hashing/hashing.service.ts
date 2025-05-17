import { randomBytes, scrypt } from 'node:crypto'
import { promisify } from 'node:util'

const scryptAsync = promisify(scrypt)

export class HashingService {
  public async hashPassword(password: string): Promise<string> {
    const salt = randomBytes(16).toString('hex')
    const derivedKey = (await scryptAsync(password, salt, 64)) as Buffer
    return `${salt}:${derivedKey.toString('hex')}`
  }

  public async verifyPassword(
    password: string,
    hash: string
  ): Promise<boolean> {
    const [salt, key] = hash.split(':')
    if (!salt || !key) {
      return false
    }
    const derivedKey = (await scryptAsync(password, salt, 64)) as Buffer
    return derivedKey.toString('hex') === key
  }
}
