import { Injectable, Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { PassportStrategy } from "@nestjs/passport";
import { ExtractJwt, Strategy } from "passport-jwt";

import { UsersService } from "../../users/users.service";
import { CurrentUserDto } from "../decorators/current-user.decorator";

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  private readonly logger: Logger = new Logger("JWT Strategy");

  constructor(
    private configService: ConfigService,
    private usersService: UsersService
  ) {
    super({
      ignoreExpiration: false,
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: configService.get<string>("JWT_API_TOKEN_SECRET"),
    });
  }

  async validate(payload: any): Promise<CurrentUserDto> {
    if (!payload.sub) {
      return null;
    }
    const { sub: id } = payload;
    return this.usersService.findOne(id);
  }
}
