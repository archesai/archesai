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
export class OrganizationRoleGuard implements CanActivate {
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

    const membership = currentUser.memberships.find(
      (val) => val.orgname == orgname
    );
    if (!membership) {
      throw new ForbiddenException(
        "You aress not authorized to access this endpoint"
      );
    }

    const roles = this.reflector.getAllAndOverride<string[]>("roles", [
      context.getHandler(),
      context.getClass(),
    ]);
    if (!roles) {
      return true;
    }

    if (!roles.includes(membership.role)) {
      throw new ForbiddenException(
        "You are not authorized to access this endpoint"
      );
    }

    return true;
  }
}
