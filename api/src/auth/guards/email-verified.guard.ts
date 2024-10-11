import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
} from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { Observable } from "rxjs";

import { CurrentUserDto } from "../decorators/current-user.decorator";

@Injectable()
export class EmailVerifiedGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  canActivate(
    context: ExecutionContext
  ): boolean | Observable<boolean> | Promise<boolean> {
    const isPublic = this.reflector.getAllAndOverride<boolean>("public", [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) {
      return true;
    }

    const { params, user } = context.switchToHttp().getRequest() as any;
    const currentUser = user as CurrentUserDto;
    const orgname = params.orgname;

    if (!orgname) {
      return true;
    }

    if (!currentUser.emailVerified) {
      throw new ForbiddenException(
        "You must verify your e-mail before using this feature."
      );
    }

    return true;
  }
}
