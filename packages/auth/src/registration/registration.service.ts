import type { UserEntity } from '@archesai/domain'

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
      deactivated: false,
      email,
      orgname:
        email.split('@')[0]! + '-' + Math.random().toString(36).substring(2, 6)
    })

    const account = await this.accountsService.create({
      authType: 'email',
      hashed_password: password,
      provider: 'LOCAL',
      providerAccountId: email,
      userId: user.id
    })

    return this.usersService.findOne(account.userId)
  }
}
