import type { ArchesApiRequest, Controller } from '@archesai/core'
import type { UserEntity } from '@archesai/domain'

import { BaseController } from '@archesai/core'
import { USER_ENTITY_KEY, UserEntitySchema } from '@archesai/domain'

import type { UsersService } from '#users/users.service'

import { CreateUserRequestSchema } from '#users/dto/create-user.req.dto'
import { UpdateUserRequestSchema } from '#users/dto/update-user.req.dto'

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
      CreateUserRequestSchema,
      UpdateUserRequestSchema,
      usersService
    )
  }

  public currentUser(request: ArchesApiRequest) {
    return request.user!
  }
}
