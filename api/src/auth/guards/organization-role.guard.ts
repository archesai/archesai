import { UserEntity } from "@/src/users/entities/user.entity";
import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
  Logger,
  NotFoundException,
} from "@nestjs/common";
import { Reflector } from "@nestjs/core";
import { Observable } from "rxjs";

@Injectable()
export class OrganizationRoleGuard implements CanActivate {
  private readonly logger = new Logger(OrganizationRoleGuard.name);
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
    const currentUser = user as UserEntity;
    const orgname = params.orgname;

    if (!orgname) {
      return true;
    }

    const membership = currentUser.memberships.find(
      (val) => val.orgname == orgname
    );
    if (!membership) {
      this.logger.error(
        `User ${currentUser.username} is not a member of organization ${orgname}`
      );
      throw new NotFoundException();
    }

    const roles = this.reflector.getAllAndOverride<string[]>("roles", [
      context.getHandler(),
      context.getClass(),
    ]);
    if (!roles) {
      return true;
    }

    if (!roles.includes(membership.role)) {
      this.logger.error(
        `User ${currentUser.username} does not have the required role in organization ${orgname}`
      );
      throw new ForbiddenException(
        "You are not authorized to access this endpoint"
      );
    }

    return true;
  }
}
