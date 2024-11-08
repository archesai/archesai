import { Injectable } from "@nestjs/common";
import { Logger } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { PassportStrategy } from "@nestjs/passport";
import { ExtractJwt, Strategy } from "passport-jwt";

import { ApiTokensService } from "../../api-tokens/api-tokens.service";
import { UsersService } from "../../users/users.service";
import { CurrentUserDto } from "../decorators/current-user.decorator";

@Injectable()
export class ApiKeyStrategy extends PassportStrategy(Strategy, "api-key-auth") {
  private readonly logger: Logger = new Logger("Api Key Strategy");

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
    this.logger.log(`Validating API Key: ${payload.id}`);
    const { id, orgname, role, uid } = payload;
    const user = await this.usersService.findOne(null, uid);
    user.memberships = user.memberships.filter((m) => m.orgname == orgname);
    if (!user.memberships.length) {
      return null;
    }

    const tokens = await this.apiTokensService.findAll(orgname, {});
    if (!tokens || !tokens.results.find((t) => t.id == id)) {
      return null;
    }
    user.memberships[0].role = role;
    return user;
  }
}
