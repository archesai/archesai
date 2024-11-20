import { UserEntity } from "@/src/users/entities/user.entity";
import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
} from "@nestjs/common";

@Injectable()
export class DeactivatedGuard implements CanActivate {
  constructor() {}
  canActivate(context: ExecutionContext) {
    const { user } = context.switchToHttp().getRequest() as any;

    // Check for user
    const currentUser = user as UserEntity;
    if (!currentUser) {
      return true;
    }

    // Check if user is deactivated
    if (currentUser.deactivated === true) {
      throw new ForbiddenException(
        "Your account has been deactivated. Please contact support."
      );
    }

    return true;
  }
}
