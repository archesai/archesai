import { UserEntity } from "@/src/users/entities/user.entity";
import {
  CanActivate,
  ExecutionContext,
  ForbiddenException,
  Injectable,
} from "@nestjs/common";

@Injectable()
export class EmailVerifiedGuard implements CanActivate {
  constructor() {}
  canActivate(context: ExecutionContext) {
    const { user } = context.switchToHttp().getRequest() as any;

    // Check for user
    const currentUser = user as UserEntity;
    if (!currentUser) {
      return true;
    }

    // Check if user is email verified
    if (currentUser.emailVerified === false) {
      throw new ForbiddenException(
        "You must verify your e-mail before using this feature."
      );
    }

    return true;
  }
}
