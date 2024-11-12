import { UserEntity } from "@/src/users/entities/user.entity";
import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
} from "@nestjs/common";
import { Reflector } from "@nestjs/core";

@Injectable()
export class DeactivatedGuard implements CanActivate {
  constructor(private reflector: Reflector) {}
  canActivate(context: ExecutionContext) {
    const isPublic = this.reflector.getAllAndOverride<boolean>("public", [
      context.getHandler(),
      context.getClass(),
    ]);
    if (isPublic) {
      return true;
    }

    // Check if deactivated
    const { user } = context.switchToHttp().getRequest() as any;
    const currentUser = user as UserEntity;
    if (currentUser?.deactivated) {
      throw new ForbiddenException();
    }

    return true;
  }
}
