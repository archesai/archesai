import { CanActivate, ExecutionContext, Injectable } from '@nestjs/common'
import { Reflector } from '@nestjs/core'
import { AuthGuard } from '@nestjs/passport'

@Injectable()
export class AppAuthGuard extends AuthGuard(['jwt', 'api-key-auth']) implements CanActivate {
  constructor(private reflector: Reflector) {
    super()
  }
  canActivate(context: ExecutionContext) {
    const isPublic = this.reflector.getAllAndOverride<boolean>('public', [context.getHandler(), context.getClass()])
    if (isPublic) {
      return true
    }

    return super.canActivate(context)
  }
}
