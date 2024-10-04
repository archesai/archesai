import { Injectable, UnauthorizedException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { PassportStrategy } from "@nestjs/passport";
import { ExtractJwt, Strategy } from "passport-jwt";

import { ApiTokensService } from "../../api-tokens/api-tokens.service";
import { UsersService } from "../../users/users.service";
import { CurrentUserDto } from "../decorators/current-user.decorator";

@Injectable()
export class ApiKeyStrategy extends PassportStrategy(Strategy, "api-key-auth") {
  constructor(
    private configService: ConfigService,
    private usersService: UsersService,
    private apiTokensService: ApiTokensService
  ) {
    super({
      ignoreExpiration: false,
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      secretOrKey: configService.get("JWT_API_TOKEN_SECRET"),
    });
  }

  async validate(payload: any): Promise<CurrentUserDto> {
    const { id, orgname, role } = payload;
    const user = await this.usersService.findOne(id);
    user.memberships = user.memberships.filter((m) => m.orgname == orgname);
    if (!user.memberships.length) {
      throw new UnauthorizedException();
    }

    const tokens = await this.apiTokensService.findAll(orgname, {});
    if (!tokens || !tokens.results.find((t) => t.id == id)) {
      throw new UnauthorizedException();
    }
    user.memberships[0].role = role;
    return user;
  }
}
