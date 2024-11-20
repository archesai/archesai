import { UserEntity } from "@/src/users/entities/user.entity";
import { Injectable, Logger } from "@nestjs/common";
import { UnauthorizedException } from "@nestjs/common";
import { PassportStrategy } from "@nestjs/passport";
import * as bcrypt from "bcryptjs";
import { Strategy } from "passport-local";

import { UsersService } from "../../users/users.service";

@Injectable()
export class LocalStrategy extends PassportStrategy(Strategy) {
  private readonly logger: Logger = new Logger("Local Strategy");

  constructor(private usersService: UsersService) {
    super({ usernameField: "email" });
  }

  async validate(email: string, password: string): Promise<UserEntity> {
    try {
      const user = await this.usersService.findOneByEmail(email);
      if (
        user &&
        user.password &&
        (await bcrypt.compare(password, user.password))
      ) {
        return user;
      }
    } catch (e) {
      this.logger.log(`Could not validate email and password: ${email}: ${e}`);
      throw new UnauthorizedException("Invalid credentials");
    }
  }
}
