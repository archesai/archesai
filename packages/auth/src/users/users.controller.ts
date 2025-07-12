import type { Controller } from '@archesai/core'
import type { UserEntity } from '@archesai/schemas'

import { BaseController } from '@archesai/core'
import {
  CreateUserDtoSchema,
  UpdateUserDtoSchema,
  USER_ENTITY_KEY,
  UserEntitySchema
} from '@archesai/schemas'

import type { UsersService } from '#users/users.service'

/**
 * Controller for handling users.
 */
export class UsersController
  extends BaseController<UserEntity>
  implements Controller
{
  constructor(usersService: UsersService) {
    super(
      USER_ENTITY_KEY,
      UserEntitySchema,
      CreateUserDtoSchema,
      UpdateUserDtoSchema,
      usersService
    )
  }
}
