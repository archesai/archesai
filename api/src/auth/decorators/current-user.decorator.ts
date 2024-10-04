import { createParamDecorator, ExecutionContext } from "@nestjs/common";
import { AuthProvider, Member, User } from "@prisma/client";

export const CurrentUser = createParamDecorator(
  (data: unknown, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest();
    return request.user as CurrentUserDto;
  }
);

export interface CurrentUserDto extends User {
  authProviders: AuthProvider[];
  memberships: Member[];
}
