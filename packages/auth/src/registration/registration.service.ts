import type { UserEntity } from '@archesai/schemas'

import { ConflictException } from '@archesai/core'

import type { AccountsService } from '#accounts/accounts.service'
import type { UsersService } from '#users/users.service'

/**
 * Service for managing Registration.
 */
export class RegistrationService {
  private readonly accountsService: AccountsService
  private readonly usersService: UsersService
  constructor(accountsService: AccountsService, usersService: UsersService) {
    this.accountsService = accountsService
    this.usersService = usersService
  }

  public async register(email: string, password: string): Promise<UserEntity> {
    // Check if the user already exists
    const exists = await this.usersService.checkIfEmailExists(email)
    if (exists) {
      throw new ConflictException('User with this email already exists')
    }

    // Create the user account
    const user = await this.usersService.create({
      createdAt: new Date().toISOString(),
      deactivated: false,
      email,
      name: email.split('@')[0] ?? email,
      // orgname:
      //   email.split('@')[0]! + '-' + Math.random().toString(36).substring(2, 6),
      updatedAt: new Date().toISOString()
    })

    const account = await this.accountsService.create({
      accountId: email,
      createdAt: new Date().toISOString(),
      password: password,
      providerId: 'LOCAL',
      updatedAt: new Date().toISOString(),
      userId: user.id
    })

    return this.usersService.findOne(account.userId)
  }
}
