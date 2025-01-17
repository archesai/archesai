import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common'
import { AuthGuard } from '@nestjs/passport'

@Injectable()
export class AuthenticatedGuard
  extends AuthGuard(['jwt', 'api-key-auth'])
  implements CanActivate
{
  constructor() {
    super()
  }
  canActivate(context: ExecutionContext) {
    return super.canActivate(context)
  }
}
