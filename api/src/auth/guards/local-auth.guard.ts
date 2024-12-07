import { ExecutionContext, Injectable } from '@nestjs/common'
import { AuthGuard } from '@nestjs/passport'

@Injectable()
export class LocalAuthGuard extends AuthGuard('local') {
  async canActivate(context: ExecutionContext) {
    return super.canActivate(context) as boolean
    // const request = context.switchToHttp().getRequest();
    // await super.logIn(request);
  }
}

// @Injectable()
// export class CookieGuard implements CanActivate {
//   async canActivate(context: ExecutionContext) {
//     const request = context.switchToHttp().getRequest();
//     return request.isAuthenticated();
//   }
// }
